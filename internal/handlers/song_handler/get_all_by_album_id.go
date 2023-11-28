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

// getByAlbumIdResponseItem represents a single song item in the GetSongsByAlbumId API response.
type getByAlbumIdResponseItem struct {
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

// getByAlbumIdResponse represents the response model for GetSongsByAlbumId API.
type getByAlbumIdResponse struct {
	// Array of songs belonging to a specific album.
	Songs []getByAlbumIdResponseItem `json:"songs"`
}

// GetByAlbumId retrieves a list of songs associated with a specific album.
// @Summary Retrieve songs by album ID
// @Description Retrieves all songs that are part of the specified album, including detailed information about each song.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   albumId   path   int     true   "Unique identifier of the album"
// @Success 200 {object} getByAlbumIdResponse "Successful response with a list of songs belonging to the requested album"
// @Failure 400 {object} response.Error "Invalid albumId format"
// @Failure 404 {object} response.Error "Album not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /albums/{albumId}/songs [get]
func (h *Handler) GetByAlbumId(c *gin.Context) {
	log.Debug().Msg("Getting songs by album")

	albumIdStr := c.Param("albumId")
	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		log.Error().Err(err).Str("albumIdStr", albumIdStr).Msg("Invalid albumId format")
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
