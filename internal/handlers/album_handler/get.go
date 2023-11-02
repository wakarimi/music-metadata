package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"modernc.org/mathutil"
	"music-metadata/internal/errors"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

// getResponse represents the response model for GetAlbum API.
type getResponse struct {
	// Unique identifier for the album.
	AlbumId int `json:"albumId"`
	// Title of the album.
	Title string `json:"title"`
	// Optional array of best cover IDs.
	BestCover *[]int `json:"bestCovers,omitempty"`
}

// Get retrieves detailed information about an album.
// @Summary Retrieve album details
// @Description Retrieves detailed information about an album, including its best covers if requested.
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param   albumId      path    int     true        "Album ID"
// @Param   bestCovers   query   int     false       "Number of best covers to retrieve"
// @Success 200 {object} getResponse
// @Failure 400 {object} response.Error "Invalid albumId or bestCovers format"
// @Failure 404 {object} response.Error "Album not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /albums/{albumId} [get]
func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting album")

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

	var album model.Album
	var bestCovers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		album, err = h.AlbumService.Get(tx, albumId)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			bestCovers, err = h.CoverService.CalcBestCoversForAlbum(tx, album.AlbumId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get album")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Album not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get album",
				Reason:  err.Error(),
			})
		}
		return
	}

	var bestCoverForAlbumResponse *[]int
	if len(bestCovers) > 0 {
		bestCoversForAlbumResponse := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers)); j++ {
				bestCoversForAlbumResponse = append(bestCoversForAlbumResponse, bestCovers[j])
			}
		}
	} else {
		bestCoverForAlbumResponse = nil
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, getResponse{
		AlbumId:   album.AlbumId,
		Title:     album.Title,
		BestCover: bestCoverForAlbumResponse,
	})
}
