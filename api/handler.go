package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	s "strings"

	"plexcache/models"
	"plexcache/plex"
	redisH "plexcache/redis"
	"plexcache/utils"

	"github.com/LukeHagar/plexgo"
	"github.com/redis/go-redis/v9"
)

func isCacheableEvent(payload models.Payload) bool {
	log.Println("payload event", payload.Event)
	return payload.Event == "media.resume" || payload.Event == "media.play"
}

func isShow(payload models.Payload) bool {
	return payload.Metadata.LibrarySectionType == "show"
}

// skipping first episode of season 1 to be sure user likes serie
func isCacheableEpisode(payload models.Payload) bool {
	return payload.Metadata.ParentIndex != 1 && payload.Metadata.Index != 1
}

func canCache(payload models.Payload) bool {
	if !isCacheableEvent(payload) || !isShow(payload) || !isCacheableEpisode(payload) {
		return false
	}

	return true
}

func isAlreadyCached(rdb *redis.Client, payload models.Payload) bool {
	ctx := context.Background()
	storedValue, err := rdb.Get(ctx, payload.Metadata.RatingKey).Result()

	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Println("Error retriving from redis", err)
		return false
	}

	var episodeCache models.EpisodeCache
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

func parsemodels(r *http.Request) (models.Payload, error) {
	var payload models.Payload
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

func copyEpisodes(episodesToCache []models.EpisodeCache) error {
	destination := "/cache"

	for _, item := range episodesToCache {
		log.Print("copy: ", item.EpisodeFilePath)
		log.Print("to: ", destination+item.EpisodeFilePath)

		err := utils.CopyFile(item.EpisodeFilePath, destination+item.EpisodeFilePath)

		if err != nil {
			return fmt.Errorf("failed to copy %s", item.Title)
		}

		for _, srtPath := range item.SrtFilePaths {
			log.Print("copy srt: ", srtPath)
			log.Print("to: ", destination+srtPath)
			err := utils.CopyFile(srtPath, destination+srtPath)

			if err != nil {
				return fmt.Errorf("failed to copy %s", srtPath)
			}
		}

	}

	return nil
}

func getSrtPaths(episodePath string, container string, stream []models.StreamPart) []string {
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

func getEpisodeCache(payload models.Payload, seasonMetadata models.SeasonMetadataResponse) []models.EpisodeCache {
	startIndex := payload.Metadata.Index + 1
	endIndex := payload.Metadata.Index + 4

	var episodesToCache []models.EpisodeCache
	for _, item := range seasonMetadata.MediaContainer.Metadata {
		if item.Index >= startIndex && item.Index <= endIndex {

			tmp := models.EpisodeCache{
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

func WebhookHandler(rdb *redis.Client, plexApi *plexgo.PlexAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Hook request")

		payload, err := parsemodels(r)

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

		seasonMetadata, err := plex.GetSeasonMetadata(plexApi, payload)

		if err != nil {
			log.Println("Failed to parse full episode response", err)
			http.Error(w, "Failed to parse full episode response", http.StatusBadRequest)
			return
		}

		episodesToCache := getEpisodeCache(payload, seasonMetadata)
		err = redisH.SaveEpisodeCacheToRedis(rdb, episodesToCache)

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
	}
}
