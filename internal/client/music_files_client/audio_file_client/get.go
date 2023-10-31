package audio_file_client

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

type GetResponse struct {
	AudioFileId       int       `json:"audioFileId"`
	Sha256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

func (c *Client) Get(audioFileId int) (audioFile GetResponse, err error) {
	log.Debug().Msg("Fetching info about audio files")

	resp, err := c.audioFileClient.Request(http.MethodGet, fmt.Sprintf("/api/audio-files/%d", audioFileId), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request for getting audio files")
		return GetResponse{}, err
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
		return GetResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return GetResponse{}, err
	}

	err = json.Unmarshal(body, &audioFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize response body")
		return GetResponse{}, err
	}

	log.Debug().Msg("Info about audio files fetched successfully")
	return audioFile, err
}
