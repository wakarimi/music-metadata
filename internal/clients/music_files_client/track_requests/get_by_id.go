package track_requests

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

type GetByIdResponse struct {
	TrackId      int    `json:"trackId"`
	CoverId      *int   `json:"coverId,omitempty"`
	DurationMs   int64  `json:"durationMs"`
	SizeByte     int64  `json:"sizeByte"`
	AudioCodec   string `json:"audioCodec"`
	BitrateKbps  int    `json:"bitrateKbps"`
	SampleRateHz int    `json:"sampleRateHz"`
	Channels     int    `json:"channels"`
	HashSha256   string `json:"hashSha256"`
}

func (c *TrackClient) GetById(trackId int) (track GetByIdResponse, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track by id from music-files service")

	resp, err := c.musicFilesClient.Request(http.MethodGet, "/api/music-files-service/tracks/"+strconv.Itoa(trackId), nil)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to execute request for getting track by ID")
		return GetByIdResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMessage := fmt.Sprintf("Received unexpected status code: %d", resp.StatusCode)
		log.Error().Err(err).Int("trackId", trackId).Str("statusCode", resp.Status).Msg(errMessage)
		return GetByIdResponse{}, fmt.Errorf(errMessage)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to read response body for track")
		return GetByIdResponse{}, err
	}

	err = json.Unmarshal(body, &track)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to deserialize response body for track")
		return GetByIdResponse{}, err
	}

	log.Debug().Int("trackId", track.TrackId).Str("sha256", track.HashSha256).Msg("Track fetched successfully")
	return track, nil
}
