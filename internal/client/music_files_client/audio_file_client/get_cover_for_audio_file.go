package audio_file_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type CoverTopForAudioFilesRequest struct {
	AudioFiles []int `json:"audioFiles"`
}

type CoverTopForAudioFilesResponse struct {
	CoversTop []int `json:"coversTop"`
}

func (c *Client) CoverTopForAudioFiles(audioFileIds []int) (coversTop []int, err error) {
	log.Debug().Msg("Fetching cover top for audio files")

	requestBody := CoverTopForAudioFilesRequest{
		AudioFiles: audioFileIds,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal requestBody")
		return make([]int, 0), err
	}

	resp, err := c.audioFileClient.Request(http.MethodGet, "/api/audio-files/covers-top", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request for fetching top for audio files")
		return make([]int, 0), err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("received unexpected status code: %d", resp.StatusCode)
		log.Error().Err(err).Str("statusCode", resp.Status).Msg("Received unexpected status code")
		return make([]int, 0), err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return make([]int, 0), err
	}

	var response CoverTopForAudioFilesResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize response body")
		return make([]int, 0), err
	}

	log.Debug().Msg("Top for audio files fetched successfully")
	return response.CoversTop, err
}
