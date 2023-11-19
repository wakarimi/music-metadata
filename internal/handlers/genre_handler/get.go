package genre_handler

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

// getResponse represents the response model for GetGenre API.
type getResponse struct {
	// Unique identifier for the genre.
	GenreId int `json:"genreId"`
	// Name of the genre.
	Name string `json:"name"`
}

// Get retrieves detailed information about a genre.
// @Summary Retrieve genre details
// @Description Retrieves detailed information about a genre, including its best covers if requested.
// @Tags Genres
// @Accept  json
// @Produce  json
// @Param   genreId      path    int     true        "Genre ID"
// @Param   bestCovers   query   int     false       "Number of best covers to retrieve"
// @Success 200 {object} getResponse
// @Failure 400 {object} response.Error "Invalid genreId or bestCovers format"
// @Failure 404 {object} response.Error "Genre not found"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /genres/{genreId} [get]
func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting genre")

	genreIdStr := c.Param("genreId")
	genreId, err := strconv.Atoi(genreIdStr)
	if err != nil {
		log.Error().Err(err).Str("genreIdStr", genreIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid genreId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("genreId", genreId).Msg("Url parameter read successfully")

	var genre model.Genre
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genre, err = h.GenreService.Get(tx, genreId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genre")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Genre not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get genre",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Genres got successfully")
	c.JSON(http.StatusOK, getResponse{
		GenreId: genre.GenreId,
		Name:    genre.Name,
	})
}
