package trackrequests

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

type GetTrackByIdResponse struct {
	TrackId   int    `json:"trackId"`
	CoverId   *int   `json:"coverId"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
}

func (c *TrackClient) GetTrackById(trackId int) (track GetTrackByIdResponse, err error) {
	resp, err := c.musicFilesClient.Request(http.MethodGet, "/api/music-files-service/tracks/"+strconv.Itoa(trackId), nil)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to execute request for getting track by ID")
		return GetTrackByIdResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMessage := fmt.Sprintf("Received unexpected status code: %d", resp.StatusCode)
		log.Error().Err(err).Int("trackId", trackId).Str("statusCode", resp.Status).Msg(errMessage)
		return GetTrackByIdResponse{}, fmt.Errorf(errMessage)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to read response body for track")
		return GetTrackByIdResponse{}, err
	}

	err = json.Unmarshal(body, &track)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to deserialize response body for track")
		return GetTrackByIdResponse{}, err
	}

	return track, nil
}