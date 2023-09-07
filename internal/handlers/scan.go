package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"music-metadata/internal/config"
	"music-metadata/internal/handlers/types"
	"net/http"
)

type ScanTrack struct {
	TrackId   int    `json:"trackId"`
	DirId     int    `json:"dirId"`
	CoverId   int    `json:"coverId"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Format    string `json:"format"`
	DateAdded string `json:"dateAdded"`
}

type ScanTracksContent struct {
	Tracks []ScanTrack `json:"tracks"`
}

func (h *MusicHandler) Scan(c *gin.Context, cfg *config.Configuration) {
	log.Info().Msg("Scanning library")

	resp, err := http.Get(cfg.MusicFilesAddress + "/api/music-files-service/tracks")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch tracks")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch tracks",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to read response",
		})
		return
	}

	var tracksResponse ScanTracksContent
	err = json.Unmarshal(body, &tracksResponse)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to unmarshal response",
		})
		return
	}

	for _, track := range tracksResponse.Tracks {
		_, err := downloadTrack(cfg.MusicFilesAddress, track.TrackId)
		if err != nil {
			log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to download track")
			continue
		}
		log.Info().Int("trackId", track.TrackId).Str("name", track.Name).Msg("Track downloaded")
	}

	c.Status(http.StatusOK)
}

func downloadTrack(baseURL string, trackId int) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/music-files-service/tracks/%d/download", baseURL, trackId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
