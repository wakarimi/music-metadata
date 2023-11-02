package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

type getByAlbumIdResponseItem struct {
	SongId      int     `json:"songId"`
	AudioFileId int     `json:"audioFileId"`
	Title       *string `json:"title"`
	AlbumId     *int    `json:"albumId"`
	ArtistId    *int    `json:"artistId"`
	GenreId     *int    `json:"genreId"`
	Year        *int    `json:"year"`
	SongNumber  *int    `json:"songNumber"`
	DiscNumber  *int    `json:"discNumber"`
	Lyrics      *string `json:"lyrics"`
	Sha256      string  `json:"sha256"`
}

type getByAlbumIdResponse struct {
	Songs []getByAlbumIdResponseItem `json:"songs"`
}

func (h *Handler) GetByAlbumId(c *gin.Context) {
	log.Debug().Msg("Getting songs by album")

	albumIdStr := c.Param("albumId")
	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		log.Error().Err(err).Str("albumIdStr", albumIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid albumId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("albumId", albumId).Msg("Url parameter read successfully")

	var songs []model.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.GetAllByAlbumId(tx, albumId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Int("albumId", albumId).Msg("Failed to get songs by album")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Album not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get songs by album",
				Reason:  err.Error(),
			})
		}
		return
	}

	songsResponseItems := make([]getByAlbumIdResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = getByAlbumIdResponseItem{
			SongId:      song.SongId,
			AudioFileId: song.AudioFileId,
			Title:       song.Title,
			AlbumId:     song.AlbumId,
			ArtistId:    song.ArtistId,
			GenreId:     song.GenreId,
			Year:        song.Year,
			SongNumber:  song.SongNumber,
			DiscNumber:  song.DiscNumber,
			Lyrics:      song.Lyrics,
			Sha256:      song.Sha256,
		}
	}

	log.Debug().Msg("Songs got successfully")
	c.JSON(http.StatusOK, getByAlbumIdResponse{
		Songs: songsResponseItems,
	})
}
