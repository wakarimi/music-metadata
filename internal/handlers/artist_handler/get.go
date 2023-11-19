package artist_handler

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

// getResponse represents the response model for GetArtist API.
type getResponse struct {
	// Unique identifier for the artist.
	ArtistId int `json:"artistId"`
	// Name of the artist.
	Name string `json:"name"`
}

// Get retrieves detailed information about an artist.
// @Summary Retrieve artist details
// @Description Retrieves detailed information about an artist, including its best covers if requested.
// @Tags Artists
// @Accept  json
// @Produce  json
// @Param   artistId      path    int     true        "Artist ID"
// @Param   bestCovers   query   int     false       "Number of best covers to retrieve"
// @Success 200 {object} getResponse
// @Failure 400 {object} response.Error "Invalid artistId or bestCovers format"
// @Failure 404 {object} response.Error "Artist not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /artists/{artistId} [get]
func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting artist")

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

	var artist model.Artist
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artist, err = h.ArtistService.Get(tx, artistId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artist")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Artist not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get artist",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Artists got successfully")
	c.JSON(http.StatusOK, getResponse{
		ArtistId: artist.ArtistId,
		Name:     artist.Name,
	})
}
