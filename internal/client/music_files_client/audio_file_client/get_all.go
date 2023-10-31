package audio_file_client

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

type GetAllResponseItem struct {
	AudioFileId       int       `json:"audioFileId"`
	Sha256            string    `json:"sha256"`
	LastContentUpdate time.Time `json:"lastContentUpdate"`
}

type GetAllResponse struct {
	AudioFiles []GetAllResponseItem `json:"audioFiles"`
}

func (c *Client) GetAll() (audioFiles []GetAllResponseItem, err error) {
	log.Debug().Msg("Fetching info about all audio files")

	resp, err := c.audioFileClient.Request(http.MethodGet, "/api/audio-files", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request for getting all audio files")
		return make([]GetAllResponseItem, 0), err
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
		return make([]GetAllResponseItem, 0), err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return make([]GetAllResponseItem, 0), err
	}

	var audioFilesObject GetAllResponse
	err = json.Unmarshal(body, &audioFilesObject)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize response body")
		return make([]GetAllResponseItem, 0), err
	}
	audioFiles = audioFilesObject.AudioFiles

	log.Debug().Int("countOfAudioFiles", len(audioFiles)).Msg("Info about all audio files fetched successfully")
	return audioFiles, err
}
