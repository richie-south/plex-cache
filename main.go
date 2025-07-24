package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Payload struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata struct {
		LibrarySectionType    string  `json:"librarySectionType"`
		RatingKey             string  `json:"ratingKey"`
		Key                   string  `json:"key"`
		ParentRatingKey       string  `json:"parentRatingKey"`
		GrandparentRatingKey  string  `json:"grandparentRatingKey"`
		GUID                  string  `json:"guid"`
		ParentGUID            string  `json:"parentGuid"`
		GrandparentGUID       string  `json:"grandparentGuid"`
		GrandparentSlug       string  `json:"grandparentSlug"`
		Type                  string  `json:"type"`
		Title                 string  `json:"title"`
		TitleSort             string  `json:"titleSort"`
		GrandparentKey        string  `json:"grandparentKey"`
		ParentKey             string  `json:"parentKey"`
		LibrarySectionTitle   string  `json:"librarySectionTitle"`
		LibrarySectionID      int     `json:"librarySectionID"`
		LibrarySectionKey     string  `json:"librarySectionKey"`
		GrandparentTitle      string  `json:"grandparentTitle"`
		ParentTitle           string  `json:"parentTitle"`
		ContentRating         string  `json:"contentRating"`
		Summary               string  `json:"summary"`
		Index                 int     `json:"index"`
		ParentIndex           int     `json:"parentIndex"`
		AudienceRating        float64 `json:"audienceRating"`
		ViewOffset            int64   `json:"viewOffset"`
		LastViewedAt          int64   `json:"lastViewedAt"`
		Year                  int     `json:"year"`
		Thumb                 string  `json:"thumb"`
		Art                   string  `json:"art"`
		ParentThumb           string  `json:"parentThumb"`
		GrandparentThumb      string  `json:"grandparentThumb"`
		GrandparentArt        string  `json:"grandparentArt"`
		GrandparentTheme      string  `json:"grandparentTheme"`
		Duration              int     `json:"duration"`
		OriginallyAvailableAt string  `json:"originallyAvailableAt"`
		AddedAt               int64   `json:"addedAt"`
		UpdatedAt             int64   `json:"updatedAt"`
		AudienceRatingImage   string  `json:"audienceRatingImage"`
		Image                 []struct {
			Alt  string `json:"alt"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"Image"`
		UltraBlurColors struct {
			TopLeft     string `json:"topLeft"`
			TopRight    string `json:"topRight"`
			BottomRight string `json:"bottomRight"`
			BottomLeft  string `json:"bottomLeft"`
		} `json:"UltraBlurColors"`
		Guid []struct {
			ID string `json:"id"`
		} `json:"Guid"`
		Rating []struct {
			Image string  `json:"image"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"Rating"`
		Director []struct {
			ID     int    `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			TagKey string `json:"tagKey"`
		} `json:"Director"`
		Writer []struct {
			ID     int     `json:"id"`
			Filter string  `json:"filter"`
			Tag    string  `json:"tag"`
			TagKey string  `json:"tagKey"`
			Thumb  *string `json:"thumb"`
		} `json:"Writer"`
		Role []struct {
			ID     int     `json:"id"`
			Filter string  `json:"filter"`
			Tag    string  `json:"tag"`
			TagKey string  `json:"tagKey"`
			Role   string  `json:"role"`
			Thumb  *string `json:"thumb"`
		} `json:"Role"`
		Producer []struct {
			ID     int    `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			TagKey string `json:"tagKey"`
		} `json:"Producer"`
	} `json:"Metadata"`
}

type PlexSearchResponse struct {
	MediaContainer struct {
		Size int `json:"size"`
		Hub  []struct {
			Title         string `json:"title"`
			Type          string `json:"type"`
			HubIdentifier string `json:"hubIdentifier"`
			Context       string `json:"context"`
			Size          int    `json:"size"`
			More          bool   `json:"more"`
			Style         string `json:"style"`
			Metadata      []struct {
				Title                 string  `json:"title"`
				Type                  string  `json:"type"`
				GrandparentTitle      string  `json:"grandparentTitle"`
				ParentTitle           string  `json:"parentTitle"`
				Summary               string  `json:"summary"`
				Year                  int     `json:"year"`
				Duration              int     `json:"duration"`
				AudienceRating        float64 `json:"audienceRating"`
				ViewOffset            int     `json:"viewOffset"`
				Thumb                 string  `json:"thumb"`
				OriginallyAvailableAt string  `json:"originallyAvailableAt"`
			}
		}
	}
}

func isCacheableEvent(payload Payload) bool {
	return payload.Event == "media.resume" || payload.Event == "media.play"
}

func isShow(payload Payload) bool {
	return payload.Metadata.LibrarySectionType == "show"
}

// skipping first episode of season 1 to be sure user likes serie
func isCacheableEpisode(payload Payload) bool {
	return payload.Metadata.ParentIndex != 1 && payload.Metadata.Index != 1
}

func canCache(payload Payload) bool {
	if !isCacheableEvent(payload) || !isShow(payload) || !isCacheableEpisode(payload) {
		return false
	}

	return true
}

func isAlreadyCached(rdb *redis.Client, payload Payload) bool {
	storedValue, err := rdb.Get(ctx, payload.Metadata.RatingKey).Result()

	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Println("Error retriving from redis", err)
		return false
	}

	var episodeCache EpisodeCache
	if err := json.Unmarshal([]byte(storedValue), &episodeCache); err != nil {
		return false
	}

	// basic bypass, if this is last cached episode allow for more to be cached
	// does not handle case of multible starting points in a series
	if episodeCache.IsLast {
		return false
	}

	return true
}

func parsePayload(r *http.Request) (Payload, error) {
	var payload Payload
	err := r.ParseMultipartForm(10 << 20) // 10 MB

	if err != nil {
		return payload, err
	}

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		return payload, err
	}

	err = json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		log.Println("Error decoding JSON:", err)
	}

	return payload, nil
}

func getSeasonMetadata(s *plexgo.PlexAPI, payload Payload) (SeasonMetadataResponse, error) {
	var fullEpisodeResponse SeasonMetadataResponse

	parentRatingKey, err := strconv.ParseFloat(payload.Metadata.ParentRatingKey, 64)
	if err != nil {
		return fullEpisodeResponse, err
	}

	metadataChildren, err := s.Library.GetMetadataChildren(ctx, parentRatingKey, plexgo.String("Stream"))
	if err != nil {
		return fullEpisodeResponse, err
	}

	body, err := io.ReadAll(metadataChildren.RawResponse.Body)
	if err != nil {
		return fullEpisodeResponse, err
	}

	err = json.Unmarshal([]byte(body), &fullEpisodeResponse)
	if err != nil {
		return fullEpisodeResponse, err
	}

	return fullEpisodeResponse, nil
}

func getEpisodeCache(payload Payload, seasonMetadata SeasonMetadataResponse) []EpisodeCache {
	startIndex := payload.Metadata.Index + 1
	endIndex := payload.Metadata.Index + 4

	var episodesToCache []EpisodeCache
	for _, item := range seasonMetadata.MediaContainer.Metadata {
		if item.Index >= startIndex && item.Index <= endIndex {

			tmp := EpisodeCache{
				RatingKey:            item.RatingKey,
				ParentRatingKey:      item.ParentRatingKey,
				GrandparentRatingKey: item.GrandparentRatingKey,
				Title:                item.Title,
				Index:                item.Index,
				ParentIndex:          item.ParentIndex,
				EpisodeFilePath:      formatEpisodePath(item.Media[0].Part[0].File),
				SrtFilePaths:         getSrtPaths(formatEpisodePath(item.Media[0].Part[0].File), item.Media[0].Part[0].Container, item.Media[0].Part[0].Stream),
				IsLast:               item.Index == endIndex,
			}

			episodesToCache = append(episodesToCache, tmp)

		}
	}

	return episodesToCache
}

func saveEpisodeCacheToRedis(rdb *redis.Client, episodesToCache []EpisodeCache) error {
	pipe := rdb.Pipeline()
	ttl := 20 * 24 * time.Hour
	for _, item := range episodesToCache {
		marshaled, err := json.Marshal(item)
		if err != nil {
			fmt.Println("could not marshal episode")
		}

		pipe.Set(ctx, item.RatingKey, marshaled, -1)
		pipe.Set(ctx, fmt.Sprintf("%s%s", item.RatingKey, ":plex-expirer"), "", ttl)
	}
	_, err := pipe.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func deleteExpiredEpisodeFromCache(rdb *redis.Client) *redis.PubSub {
	plexExpirerKey := ":plex-expirer"
	location := "/cache"

	subscriber := rdb.PSubscribe(ctx, "__keyevent@0__:expired")
	defer subscriber.Close()

	for msg := range subscriber.Channel() {

		if !s.HasSuffix(msg.Payload, plexExpirerKey) {
			continue
		}

		dataKey := s.Split(msg.Payload, plexExpirerKey)[0]
		log.Println("Time to remove", dataKey)
		storedValue, err := rdb.Get(ctx, dataKey).Result()

		if err == redis.Nil {
			log.Println("redis.Nil", err)
			continue
		} else if err != nil {
			log.Println("Error retriving key to delete from redis", err)
			continue
		}

		var episodeCache EpisodeCache
		if err := json.Unmarshal([]byte(storedValue), &episodeCache); err != nil {
			continue
		}

		err = RemoveFile(location + episodeCache.EpisodeFilePath)
		log.Println("Removed", episodeCache.EpisodeFilePath)

		if err != nil {
			log.Println("failed to remove file", err)
			continue
		}

		for _, srtPath := range episodeCache.SrtFilePaths {
			err = RemoveFile(location + srtPath)
			log.Println("Removed srt", srtPath)

			if err != nil {
				log.Println("failed to remove srt", err)
				continue
			}
		}

		err = rdb.Del(ctx, dataKey).Err()

		if err != nil {
			log.Println("failed to remove dataKey from redis", err)
			continue
		}

	}

	return subscriber
}

func main() {
	log.Println("Started")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Error connecing to redis %e", err)
		return
	}

	subscriber := deleteExpiredEpisodeFromCache(rdb)
	defer subscriber.Close()

	s := plexgo.New(
		plexgo.WithSecurity(os.Getenv("PLEX_API_KEY")),
		plexgo.WithIP(os.Getenv("PLEX_IP")),
		plexgo.WithProtocol("http"),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hook request")

		payload, err := parsePayload(r)

		if err != nil {
			log.Println("could not parse payload")
			http.Error(w, "No payload found", http.StatusBadRequest)
			return
		}

		if isAlreadyCached(rdb, payload) {
			log.Println("Already cached")
			w.WriteHeader(http.StatusOK)
			return
		}

		if !canCache(payload) {
			log.Println("should not cache")
			w.WriteHeader(http.StatusOK)
			return
		}

		seasonMetadata, err := getSeasonMetadata(s, payload)

		if err != nil {
			log.Println("Failed to parse full episode response", err)
			http.Error(w, "Failed to parse full episode response", http.StatusBadRequest)
			return
		}

		episodesToCache := getEpisodeCache(payload, seasonMetadata)
		err = saveEpisodeCacheToRedis(rdb, episodesToCache)

		if err != nil {
			log.Println("could not pipe to redis")
		}

		err = copyEpisodes(episodesToCache)

		if err != nil {
			log.Println("could not move files", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("request ok")
	}).Methods("POST")

	http.ListenAndServe(":4001", r)
}

func getSrtPaths(episodePath string, container string, stream []StreamPart) []string {
	var srtFilePaths []string
	for _, item := range stream {
		if item.Format == "srt" {
			srtFilePaths = append(srtFilePaths, s.Replace(episodePath, container, item.LanguageTag+"."+item.Format, 1))
		}
	}

	return srtFilePaths
}

func formatEpisodePath(episodePath string) string {
	return s.Replace(episodePath, "/data/tvshows", "/media/tvshows", 1)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func RemoveFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path %s is a directory, not a file", path)
	}

	return os.Remove(path)
}

func copyEpisodes(episodesToCache []EpisodeCache) error {
	destination := "/cache"

	for _, item := range episodesToCache {
		log.Print("copy: ", item.EpisodeFilePath)
		log.Print("to: ", destination+item.EpisodeFilePath)

		err := CopyFile(item.EpisodeFilePath, destination+item.EpisodeFilePath)

		if err != nil {
			return fmt.Errorf("failed to copy %s", item.Title)
		}

		for _, srtPath := range item.SrtFilePaths {
			log.Print("copy srt: ", srtPath)
			log.Print("to: ", destination+srtPath)
			err := CopyFile(srtPath, destination+srtPath)

			if err != nil {
				return fmt.Errorf("failed to copy %s", srtPath)
			}
		}

	}

	return nil
}

type SeasonMetadataResponse struct {
	MediaContainer struct {
		Size                     int    `json:"size"`
		AllowSync                bool   `json:"allowSync"`
		Art                      string `json:"art"`
		GrandparentContentRating string `json:"grandparentContentRating"`
		GrandparentRatingKey     int    `json:"grandparentRatingKey"`
		GrandparentStudio        string `json:"grandparentStudio"`
		GrandparentTheme         string `json:"grandparentTheme"`
		GrandparentThumb         string `json:"grandparentThumb"`
		GrandparentTitle         string `json:"grandparentTitle"`
		Identifier               string `json:"identifier"`
		Key                      string `json:"key"`
		LibrarySectionID         int    `json:"librarySectionID"`
		LibrarySectionTitle      string `json:"librarySectionTitle"`
		LibrarySectionUUID       string `json:"librarySectionUUID"`
		MediaTagPrefix           string `json:"mediaTagPrefix"`
		MediaTagVersion          int    `json:"mediaTagVersion"`
		Nocache                  bool   `json:"nocache"`
		ParentIndex              int    `json:"parentIndex"`
		ParentTitle              string `json:"parentTitle"`
		Theme                    string `json:"theme"`
		Thumb                    string `json:"thumb"`
		Title1                   string `json:"title1"`
		Title2                   string `json:"title2"`
		ViewGroup                string `json:"viewGroup"`
		Metadata                 []struct {
			RatingKey             string  `json:"ratingKey"`
			Key                   string  `json:"key"`
			ParentRatingKey       string  `json:"parentRatingKey"`
			GrandparentRatingKey  string  `json:"grandparentRatingKey"`
			GUID                  string  `json:"guid"`
			ParentGUID            string  `json:"parentGuid"`
			GrandparentGUID       string  `json:"grandparentGuid"`
			GrandparentSlug       string  `json:"grandparentSlug"`
			Type                  string  `json:"type"`
			Title                 string  `json:"title"`
			TitleSort             string  `json:"titleSort"`
			GrandparentKey        string  `json:"grandparentKey"`
			ParentKey             string  `json:"parentKey"`
			GrandparentTitle      string  `json:"grandparentTitle"`
			ParentTitle           string  `json:"parentTitle"`
			ContentRating         string  `json:"contentRating"`
			Summary               string  `json:"summary"`
			Index                 int     `json:"index"`
			ParentIndex           int     `json:"parentIndex"`
			AudienceRating        float64 `json:"audienceRating"`
			ViewCount             int     `json:"viewCount"`
			LastViewedAt          int64   `json:"lastViewedAt"`
			Year                  int     `json:"year"`
			Thumb                 string  `json:"thumb"`
			Art                   string  `json:"art"`
			ParentThumb           string  `json:"parentThumb"`
			GrandparentThumb      string  `json:"grandparentThumb"`
			GrandparentArt        string  `json:"grandparentArt"`
			GrandparentTheme      string  `json:"grandparentTheme"`
			Duration              int64   `json:"duration"`
			OriginallyAvailableAt string  `json:"originallyAvailableAt"`
			AddedAt               int64   `json:"addedAt"`
			UpdatedAt             int64   `json:"updatedAt"`
			AudienceRatingImage   string  `json:"audienceRatingImage"`
			Media                 []struct {
				ID               int     `json:"id"`
				Duration         int64   `json:"duration"`
				Bitrate          int     `json:"bitrate"`
				Width            int     `json:"width"`
				Height           int     `json:"height"`
				AspectRatio      float64 `json:"aspectRatio"`
				AudioChannels    int     `json:"audioChannels"`
				AudioCodec       string  `json:"audioCodec"`
				VideoCodec       string  `json:"videoCodec"`
				VideoResolution  string  `json:"videoResolution"`
				Container        string  `json:"container"`
				VideoFrameRate   string  `json:"videoFrameRate"`
				AudioProfile     string  `json:"audioProfile"`
				VideoProfile     string  `json:"videoProfile"`
				HasVoiceActivity bool    `json:"hasVoiceActivity"`
				Part             []struct {
					ID           int          `json:"id"`
					Key          string       `json:"key"`
					Duration     int64        `json:"duration"`
					File         string       `json:"file"`
					Size         int64        `json:"size"`
					AudioProfile string       `json:"audioProfile"`
					Container    string       `json:"container"`
					VideoProfile string       `json:"videoProfile"`
					Stream       []StreamPart `json:"Stream"`
				}
			}
		}
	}
}

type StreamPart struct {
	ID                   int     `json:"id"`
	StreamType           int     `json:"streamType"`
	Default              bool    `json:"default"`
	Selected             bool    `json:"selected,omitempty"`
	Codec                string  `json:"codec"`
	Index                int     `json:"index"`
	Bitrate              int     `json:"bitrate,omitempty"`
	Channels             int     `json:"channels,omitempty"`
	Profile              string  `json:"profile,omitempty"`
	SamplingRate         int     `json:"samplingRate,omitempty"`
	Language             string  `json:"language"`
	LanguageTag          string  `json:"languageTag"`
	LanguageCode         string  `json:"languageCode"`
	BitDepth             int     `json:"bitDepth,omitempty"`
	ChromaLocation       string  `json:"chromaLocation,omitempty"`
	ChromaSubsampling    string  `json:"chromaSubsampling,omitempty"`
	CodedHeight          int     `json:"codedHeight,omitempty"`
	CodedWidth           int     `json:"codedWidth,omitempty"`
	ColorPrimaries       string  `json:"colorPrimaries,omitempty"`
	ColorRange           string  `json:"colorRange,omitempty"`
	ColorSpace           string  `json:"colorSpace,omitempty"`
	ColorTrc             string  `json:"colorTrc,omitempty"`
	FrameRate            float64 `json:"frameRate,omitempty"`
	HasScalingMatrix     bool    `json:"hasScalingMatrix,omitempty"`
	Height               int     `json:"height,omitempty"`
	Level                int     `json:"level,omitempty"`
	RefFrames            int     `json:"refFrames,omitempty"`
	ScanType             string  `json:"scanType,omitempty"`
	Width                int     `json:"width,omitempty"`
	DisplayTitle         string  `json:"displayTitle"`
	ExtendedDisplayTitle string  `json:"extendedDisplayTitle"`
	CanAutoSync          bool    `json:"canAutoSync,omitempty"`
	Format               string  `json:"format,omitempty"`
	Key                  string  `json:"key,omitempty"`
}

type EpisodeCache struct {
	RatingKey            string   `json:"ratingKey"`
	ParentRatingKey      string   `json:"parentRatingKey"`
	GrandparentRatingKey string   `json:"grandparentRatingKey"`
	Title                string   `json:"title"`
	Index                int      `json:"index"`
	ParentIndex          int      `json:"parentIndex"`
	EpisodeFilePath      string   `json:"episodeFilePath"`
	SrtFilePaths         []string `json:"srtFilePaths"`
	IsLast               bool     `json:"isLast"`
}
