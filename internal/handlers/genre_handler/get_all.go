package genre_handler

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

// getAllResponseItem represents a single genre item in the GetAllGenres API response.
type getAllResponseItem struct {
	// Unique identifier for the genre.
	GenreId int `json:"genreId"`
	// Name of the genre.
	Name string `json:"name"`
	// Optional array of best cover IDs for the genre.
	BestCover *[]int `json:"bestCovers,omitempty"`
}

// getAllResponse represents the response model for GetAllGenres API.
type getAllResponse struct {
	// Array of genres.
	Genres []getAllResponseItem `json:"genres"`
}

// GetAll retrieves a list of all genres with optional best covers.
// @Summary Retrieve all genres
// @Description Retrieves a list of all genres, including their best covers if requested.
// @Tags Genres
// @Accept  json
// @Produce  json
// @Param   bestCovers   query   int     false       "Number of best covers for each genre to retrieve"
// @Success 200 {object} getAllResponse "Success response with a list of genres and optional best covers for each"
// @Failure 400 {object} response.Error "Invalid bestCovers format"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /genres [get]
func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting genres")

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

	var genres []model.Genre
	var bestCovers [][]int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genres, err = h.GenreService.GetAll(tx)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			for _, genre := range genres {
				bestCoversItem, err := h.CoverService.CalcBestCoversForGenre(tx, genre.GenreId)
				if err != nil {
					return err
				}
				bestCovers = append(bestCovers, bestCoversItem)
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genres")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get genres",
			Reason:  err.Error(),
		})
		return
	}

	genresResponseItems := make([]getAllResponseItem, len(genres))
	for i, genre := range genres {
		bestCoversForGenre := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers[i])); j++ {
				bestCoversForGenre = append(bestCoversForGenre, bestCovers[i][j])
			}
		}
		var bestCoverForGenreResponse *[]int
		if len(bestCoversForGenre) > 0 {
			bestCoverForGenreResponse = &bestCoversForGenre
		} else {
			bestCoverForGenreResponse = nil
		}
		genresResponseItems[i] = getAllResponseItem{
			GenreId:   genre.GenreId,
			Name:      genre.Name,
			BestCover: bestCoverForGenreResponse,
		}
	}

	log.Debug().Msg("Genres got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Genres: genresResponseItems,
	})
}
