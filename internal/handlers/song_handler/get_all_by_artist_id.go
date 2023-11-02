package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

type getByArtistIdResponseItem struct {
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

type getByArtistIdResponse struct {
	Songs []getByArtistIdResponseItem `json:"songs"`
}

func (h *Handler) GetByArtistId(c *gin.Context) {
	log.Debug().Msg("Getting songs by artist")

	artistIdStr := c.Param("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil {
		log.Error().Err(err).Str("artistIdStr", artistIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid artistId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("artistId", artistId).Msg("Url parameter read successfully")

	var songs []model.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.GetAllByArtistId(tx, artistId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to get songs")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get songs",
			Reason:  err.Error(),
		})
		return
	}

	songsResponseItems := make([]getByArtistIdResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = getByArtistIdResponseItem{
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
	c.JSON(http.StatusOK, getByArtistIdResponse{
		Songs: songsResponseItems,
	})
}
