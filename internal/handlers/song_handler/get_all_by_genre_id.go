package song_handler

import (
	"music-metadata/internal/errors"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// getByGenreIdResponseItem represents a single song item in the GetSongsByGenreId API response.
type getByGenreIdResponseItem struct {
	// Unique identifier for the song.
	SongId int `json:"songId"`
	// Identifier for the associated audio file.
	AudioFileId int `json:"audioFileId"`
	// Title of the song.
	Title *string `json:"title"`
	// Identifier of the album to which the song belongs.
	AlbumId *int `json:"albumId"`
	// Identifier of the artist of the song.
	ArtistId *int `json:"artistId"`
	// Genre identifier of the song.
	GenreId *int `json:"genreId"`
	// Release year of the song.
	Year *int `json:"year"`
	// Track number of the song in the album.
	SongNumber *int `json:"songNumber"`
	// Disc number of the song in the album.
	DiscNumber *int `json:"discNumber"`
	// Lyrics of the song.
	Lyrics *string `json:"lyrics"`
	// SHA256 hash of the song file.
	Sha256 string `json:"sha256"`
}

// getByGenreIdResponse represents the response model for GetSongsByGenreId API.
type getByGenreIdResponse struct {
	// Array of songs belonging to a specific artist.
	Songs []getByGenreIdResponseItem `json:"songs"`
}

// GetByGenreId retrieves a list of songs associated with a specific genre.
// @Summary Retrieve songs by genre ID
// @Description Retrieves all songs that are part of the specified genre, including detailed information about each song.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   genreId   path   int     true   "Unique identifier of the genre"
// @Success 200 {object} getByGenreIdResponse "Successful response with a list of songs belonging to the requested genre"
// @Failure 400 {object} response.Error "Invalid genreId format"
// @Failure 404 {object} response.Error "Genre not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /genres/{genreId}/songs [get]
func (h *Handler) GetByGenreId(c *gin.Context) {
	log.Debug().Msg("Getting songs by genre")

	genreIdStr := c.Param("genreId")
	genreId, err := strconv.Atoi(genreIdStr)
	if err != nil {
		log.Error().Err(err).Str("genreIdStr", genreIdStr).Msg("Invalid genreId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid genreId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("genreId", genreId).Msg("Url parameter read successfully")

	var songs []model.Song
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		songs, err = h.SongService.GetAllByGenreId(tx, genreId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Int("genreId", genreId).Msg("Failed to get songs by genre")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Album not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get songs by genre",
				Reason:  err.Error(),
			})
		}
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
