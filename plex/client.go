package plex

import (
	"context"
	"encoding/json"
	"io"
	"plexcache/models"
	"strconv"

	"github.com/LukeHagar/plexgo"
)

func GetSeasonMetadata(s *plexgo.PlexAPI, payload models.Payload) (models.SeasonMetadataResponse, error) {
	ctx := context.Background()
	var fullEpisodeResponse models.SeasonMetadataResponse

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
