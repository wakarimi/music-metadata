package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"music-metadata/internal/config"
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

func Scan(c *gin.Context, cfg *config.Configuration) {
	log.Println("Что")
	resp, err := http.Get(cfg.MusicFilesAddress + "/api/music-files-service/tracks")
	if err != nil {
		log.Println("Failed to fetch tracks:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tracks"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var tracksResponse ScanTracksContent
	err = json.Unmarshal(body, &tracksResponse)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	log.Println(tracksResponse)
}
