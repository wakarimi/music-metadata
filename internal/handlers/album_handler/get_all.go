package album_handler

import (
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// getAllResponseItem represents a single album item in the GetAllAlbums API response.
type getAllResponseItem struct {
	// Unique identifier for the album.
	AlbumId int `json:"albumId"`
	// Title of the album.
	Title string `json:"title"`
}

// GetAll retrieves a list of all albums with optional best covers.
// @Summary Retrieve all albums
// @Description Retrieves a list of all albums, including their best covers if requested.
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param   bestCovers   query   int     false       "Number of best covers for each album to retrieve"
// @Success 200 {object} getAllResponse "Success response with a list of albums and optional best covers for each"
// @Failure 400 {object} response.Error "Invalid bestCovers format"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /albums [get]
func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting albums")

	var albums []model.Album
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albums, err = h.AlbumService.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get albums")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get albums",
			Reason:  err.Error(),
		})
		return
	}

	albumsResponseItems := make([]getAllResponseItem, len(albums))
	for i, album := range albums {
		albumsResponseItems[i] = getAllResponseItem{
			AlbumId: album.AlbumId,
			Title:   album.Title,
		}
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, albumsResponseItems)
}
