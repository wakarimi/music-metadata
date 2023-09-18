package track_metadata

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
	"net/http"
	"strings"
)

type ScanTrack struct {
	TrackId    int    `json:"trackId"`
	CoverId    *int   `json:"coverId,omitempty"`
	DurationMs int64  `json:"durationMs"`
	AudioCodec string `json:"audioCodec"`
	SizeByte   int64  `json:"sizeByte"`
	HashSha256 string `json:"hashSha256"`
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

		trackMetadataExisted, err := h.TrackRepo.IsExistsByTrackId(track.TrackId)
		if err != nil {
			log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to check track metadata existence")
			continue
		}
		if trackMetadataExisted {
			trackMetadata, err := h.TrackRepo.ReadByTrackId(track.TrackId)
			if err != nil {
				log.Error().Err(err).Int("trackId", track.TrackId).Msg("Failed to get track metadata")
				continue
			}
			err = h.TrackRepo.Update(trackMetadata.TrackMetadataId, models.TrackMetadata{
				TrackId:     track.TrackId,
				Title:       h.getTitle(metadata),
				AlbumId:     h.getOrCreateAlbum(metadata),
				ArtistId:    h.getOrCreateArtist(metadata),
				GenreId:     h.getOrCreateGenre(metadata),
				Year:        h.getYear(metadata),
				TrackNumber: h.getTrackNumber(metadata),
				DiscNumber:  h.getDiscNumber(metadata),
				Lyrics:      h.getLyrics(metadata),
				HashSha256:  track.HashSha256,
			})
			if err != nil {
				continue
			}
		} else {
			_, err = h.TrackRepo.Create(models.TrackMetadata{
				TrackId:     track.TrackId,
				Title:       h.getTitle(metadata),
				AlbumId:     h.getOrCreateAlbum(metadata),
				ArtistId:    h.getOrCreateArtist(metadata),
				GenreId:     h.getOrCreateGenre(metadata),
				Year:        h.getYear(metadata),
				TrackNumber: h.getTrackNumber(metadata),
				DiscNumber:  h.getDiscNumber(metadata),
				Lyrics:      h.getLyrics(metadata),
				HashSha256:  track.HashSha256,
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

func (h *MusicHandler) getTitle(metadata tag.Metadata) *string {
	title := strings.TrimSpace(metadata.Title())

	if len(title) == 0 {
		return nil
	}

	return &title
}

func (h *MusicHandler) getYear(metadata tag.Metadata) *int {
	year := metadata.Year()

	if year == 0 {
		return nil
	}

	return &year
}

func (h *MusicHandler) getTrackNumber(metadata tag.Metadata) *int {
	trackNumber, _ := metadata.Track()

	if trackNumber == 0 {
		return nil
	}

	return &trackNumber
}

func (h *MusicHandler) getDiscNumber(metadata tag.Metadata) *int {
	discNumber, _ := metadata.Disc()

	if discNumber == 0 {
		return nil
	}

	return &discNumber
}

func (h *MusicHandler) getLyrics(metadata tag.Metadata) *string {
	lyrics := strings.TrimSpace(metadata.Lyrics())

	if len(lyrics) == 0 {
		return nil
	}

	return &lyrics
}

func (h *MusicHandler) getOrCreateAlbum(metadata tag.Metadata) (albumId *int) {
	title := strings.TrimSpace(metadata.Album())

	if len(title) == 0 {
		return nil
	}

	albumExists, err := h.AlbumRepo.IsExistsByTitle(title)
	if err != nil {
		return nil
	}

	if albumExists {
		album, err := h.AlbumRepo.ReadByTitle(title)
		if err != nil {
			return nil
		}
		albumId = &album.AlbumId
	} else {
		newAlbumId, err := h.AlbumRepo.Create(models.Album{
			Title: title,
		})
		if err != nil {
			return nil
		}
		albumId = &newAlbumId
	}
	return albumId
}

func (h *MusicHandler) getOrCreateArtist(metadata tag.Metadata) (artistId *int) {
	name := strings.TrimSpace(metadata.Artist())

	if len(name) == 0 {
		return nil
	}

	artistExists, err := h.ArtistRepo.IsExistsByName(name)
	if err != nil {
		return nil
	}

	if artistExists {
		artist, err := h.ArtistRepo.ReadByName(name)
		if err != nil {
			return nil
		}
		artistId = &artist.ArtistId
	} else {
		newArtistId, err := h.ArtistRepo.Create(models.Artist{
			Name: name,
		})
		if err != nil {
			return nil
		}
		artistId = &newArtistId
	}
	return artistId
}

func (h *MusicHandler) getOrCreateGenre(metadata tag.Metadata) (genreId *int) {
	name := strings.TrimSpace(metadata.Genre())

	if len(name) == 0 {
		return nil
	}

	genreExists, err := h.GenreRepo.IsExistsByName(name)
	if err != nil {
		return nil
	}

	if genreExists {
		genre, err := h.GenreRepo.ReadByName(name)
		if err != nil {
			return nil
		}
		genreId = &genre.GenreId
	} else {
		newGenreId, err := h.GenreRepo.Create(models.Genre{
			Name: name,
		})
		if err != nil {
			return nil
		}
		genreId = &newGenreId
	}
	return genreId
}
