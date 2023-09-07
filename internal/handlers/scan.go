package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"music-metadata/internal/config"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/models"
	"music-metadata/internal/utils"
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

func (h *MusicHandler) Scan(c *gin.Context, httpServerConfig *config.HttpServer) {
	log.Info().Msg("Scanning library")

	resp, err := http.Get(httpServerConfig.OtherHttpServers.MusicFilesAddress + "/api/music-files-service/tracks")
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
		trackData, err := downloadTrack(httpServerConfig.MusicFilesAddress, track.TrackId)
		if err != nil {
			log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to download track")
			continue
		}

		metadata, err := extractMetadata(trackData)
		if err != nil {
			log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to extract tags from track")
			continue
		}

		var albumId *int
		albumExists, err := h.AlbumRepo.IsAlbumExistsByTitle(metadata.Album())
		if err != nil {
			continue
		}
		if albumExists {
			album, err := h.AlbumRepo.ReadAlbumByTitle(metadata.Album())
			if err != nil {
				continue
			}
			albumId = &album.AlbumId
		} else {
			albumId, err = h.AlbumRepo.CreateAlbum(models.Album{
				Title: metadata.Album(),
			})
			if err != nil {
				continue
			}
		}

		var artistId *int
		artistExists, err := h.ArtistRepo.IsArtistExistsByName(metadata.Artist())
		if err != nil {
			continue
		}
		if artistExists {
			artist, err := h.ArtistRepo.ReadArtistByName(metadata.Artist())
			if err != nil {
				continue
			}
			artistId = &artist.ArtistId
		} else {
			artistId, err = h.ArtistRepo.CreateArtist(models.Artist{
				Name: metadata.Artist(),
			})
			if err != nil {
				continue
			}
		}

		var genreId *int
		genreExists, err := h.GenreRepo.IsGenreExistsByName(metadata.Genre())
		if err != nil {
			continue
		}
		if genreExists {
			genre, err := h.GenreRepo.ReadGenreByName(metadata.Genre())
			if err != nil {
				continue
			}
			genreId = &genre.GenreId
		} else {
			genreId, err = h.GenreRepo.CreateGenre(models.Genre{
				Name: metadata.Genre(),
			})
			if err != nil {
				continue
			}
		}

		trackMetadataExisted, err := h.TrackRepo.IsTrackMetadataExistsByTrackId(track.TrackId)
		if err != nil {
			log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to check track metadata existence")
			continue
		}
		if trackMetadataExisted {
			trackMetadata, err := h.TrackRepo.ReadTrackMetadataByTrackId(track.TrackId)
			if err != nil {
				log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to get track metadata")
				continue
			}
			err = h.TrackRepo.UpdateTrackMetadata(trackMetadata.TrackMetadataId, models.TrackMetadata{
				Title:      utils.StringToPointer(metadata.Title()),
				ArtistId:   artistId,
				AlbumId:    albumId,
				Genre:      genreId,
				Bitrate:    nil,
				Channels:   nil,
				SampleRate: nil,
				Duration:   nil,
			})
			if err != nil {
				continue
			}
		} else {
			_, err = h.TrackRepo.CreateTrackMetadata(models.TrackMetadata{
				TrackId:    track.TrackId,
				Title:      utils.StringToPointer(metadata.Title()),
				ArtistId:   artistId,
				AlbumId:    albumId,
				Genre:      genreId,
				Bitrate:    nil,
				Channels:   nil,
				SampleRate: nil,
				Duration:   nil,
			})
			if err != nil {
				continue
			}
		}

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

func extractMetadata(trackData []byte) (metadata tag.Metadata, err error) {
	r := bytes.NewReader(trackData)
	return tag.ReadFrom(r)
}
