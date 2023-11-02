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

// getResponse represents a single song item in the response of the Get API.
type getResponse struct {
	// SongId is the unique identifier for the song.
	SongId int `json:"songId"`
	// AudioFileId is the identifier of the associated audio file.
	AudioFileId int `json:"audioFileId"`
	// Title is the title of the song.
	Title *string `json:"title"`
	// AlbumId is the identifier of the album to which the song belongs.
	AlbumId *int `json:"albumId"`
	// ArtistId is the identifier of the song's artist.
	ArtistId *int `json:"artistId"`
	// GenreId is the genre identifier of the song.
	GenreId *int `json:"genreId"`
	// Year is the release year of the song.
	Year *int `json:"year"`
	// SongNumber is the track number of the song in the album.
	SongNumber *int `json:"songNumber"`
	// DiscNumber is the disc number of the song in the album.
	DiscNumber *int `json:"discNumber"`
	// Lyrics are the lyrics of the song.
	Lyrics *string `json:"lyrics"`
	// Sha256 is the SHA256 hash of the song file.
	Sha256 string `json:"sha256"`
}

// Get handles the request to retrieve a specific song by its ID.
// @Summary Retrieve a song by its ID
// @Description Retrieves detailed information about a song specified by its unique ID.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   songId     path   int     true   "Unique identifier of the song"
// @Success 200 {object} getResponse "Successful response with song details"
// @Failure 400 {object} response.Error "Invalid songId format"
// @Failure 404 {object} response.Error "Song not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /songs/{songId} [get]
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
