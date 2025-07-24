package models

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
