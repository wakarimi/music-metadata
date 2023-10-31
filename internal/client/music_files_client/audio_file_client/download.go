package audio_file_client

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func (c *Client) Download(audioFileId int) (file []byte, err error) {
	log.Debug().Msg("Fetching info about all audio files")

	resp, err := c.audioFileClient.Request(http.MethodGet, fmt.Sprintf("/api/audio-files/%d/download", audioFileId), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request for download audio file")
		return make([]byte, 0), err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download audio file with id=%d", audioFileId)
		return make([]byte, 0), err
	}

	return io.ReadAll(resp.Body)
}
