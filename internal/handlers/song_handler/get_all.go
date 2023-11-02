package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
)

// getAllResponseItem represents a single song item in the GetAll API response.
type getAllResponseItem struct {
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

// getAllResponse wraps the list of songs in the GetAll API response.
type getAllResponse struct {
	// Songs is an array of song items.
	Songs []getAllResponseItem `json:"songs"`
}

// GetAll handles the request to retrieve a list of all songs.
// @Summary Retrieve a list of all songs
// @Description Retrieves detailed information about all available songs.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Success 200 {array} getAllResponseItem "Successful response with list of songs"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /songs [get]
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

	songsResponseItems := make([]getAllResponseItem, len(songs))
	for i, song := range songs {
		songsResponseItems[i] = getAllResponseItem{
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
	c.JSON(http.StatusOK, getAllResponse{
		Songs: songsResponseItems,
	})
}
