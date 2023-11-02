package genre_handler

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

// getResponse represents the response model for GetGenre API.
type getResponse struct {
	// Unique identifier for the genre.
	GenreId int `json:"genreId"`
	// Name of the genre.
	Name string `json:"name"`
	// Optional array of best cover IDs.
	BestCover *[]int `json:"bestCovers,omitempty"`
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

	var genre model.Genre
	var bestCovers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genre, err = h.GenreService.Get(tx, genreId)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			bestCovers, err = h.CoverService.CalcBestCoversForGenre(tx, genre.GenreId)
			if err != nil {
				return err
			}
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

	var bestCoverForGenreResponse *[]int
	if len(bestCovers) > 0 {
		bestCoversForGenreResponse := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers)); j++ {
				bestCoversForGenreResponse = append(bestCoversForGenreResponse, bestCovers[j])
			}
		}
	} else {
		bestCoverForGenreResponse = nil
	}

	log.Debug().Msg("Genres got successfully")
	c.JSON(http.StatusOK, getResponse{
		GenreId:   genre.GenreId,
		Name:      genre.Name,
		BestCover: bestCoverForGenreResponse,
	})
}
