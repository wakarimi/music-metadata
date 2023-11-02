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

type getResponse struct {
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

func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting song")

	songIdStr := c.Param("songId")
	songId, err := strconv.Atoi(songIdStr)
	if err != nil {
		log.Error().Err(err).Str("songIdStr", songIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid songId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("songId", songId).Msg("Url parameter read successfully")

	var song model.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		song, err = h.SongService.Get(tx, songId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get song")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Song not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get song",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Songs got successfully")
	c.JSON(http.StatusOK, getResponse{
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
	})
}
