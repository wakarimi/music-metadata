package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
)

type getAllResponseItem struct {
	SongId      int     `db:"songId"`
	AudioFileId int     `db:"audioFileId"`
	Title       *string `db:"title"`
	AlbumId     *int    `db:"albumId"`
	ArtistId    *int    `db:"artistId"`
	GenreId     *int    `db:"genreId"`
	Year        *int    `db:"year"`
	SongNumber  *int    `db:"songNumber"`
	DiscNumber  *int    `db:"discNumber"`
	Lyrics      *string `db:"lyrics"`
	Sha256      string  `db:"sha256"`
}

type getAllResponse struct {
	Songs []getByGenreIdResponseItem `json:"songs"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting songs")

	var songs []model.Song
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get songs")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get songs",
			Reason:  err.Error(),
		})
		return
	}

	songsResponseItems := make([]getByGenreIdResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = getByGenreIdResponseItem{
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
	c.JSON(http.StatusOK, getByGenreIdResponse{
		Songs: songsResponseItems,
	})
}
