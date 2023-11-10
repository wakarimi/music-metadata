package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"modernc.org/mathutil"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

// getAllResponseItem represents a single album item in the GetAllAlbums API response.
type getAllResponseItem struct {
	// Unique identifier for the album.
	AlbumId int `json:"albumId"`
	// Title of the album.
	Title string `json:"title"`
	// Optional array of best cover IDs for the album.
	BestCover *[]int `json:"bestCovers,omitempty"`
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

	bestCoversStr := c.DefaultQuery("bestCovers", "0")
	bestCoversInt, err := strconv.Atoi(bestCoversStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bestCovers format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid bestCovers format",
			Reason:  err.Error(),
		})
		return
	}

	var albums []model.Album
	var bestCovers [][]int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albums, err = h.AlbumService.GetAll(tx)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			for _, album := range albums {
				bestCoversItem, err := h.CoverService.CalcBestCoversForAlbum(tx, album.AlbumId)
				if err != nil {
					return err
				}
				bestCovers = append(bestCovers, bestCoversItem)
			}
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
		bestCoversForAlbum := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers[i])); j++ {
				bestCoversForAlbum = append(bestCoversForAlbum, bestCovers[i][j])
			}
		}
		var bestCoverForAlbumResponse *[]int
		if len(bestCoversForAlbum) > 0 {
			bestCoverForAlbumResponse = &bestCoversForAlbum
		} else {
			bestCoverForAlbumResponse = nil
		}
		albumsResponseItems[i] = getAllResponseItem{
			AlbumId:   album.AlbumId,
			Title:     album.Title,
			BestCover: bestCoverForAlbumResponse,
		}
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, albumsResponseItems)
}
