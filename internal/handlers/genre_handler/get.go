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

type getResponse struct {
	GenreId   int    `json:"genreId"`
	Name      string `json:"name"`
	BestCover *[]int `json:"bestCovers,omitempty"`
}

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

	bestCoversForGenreResponse := make([]int, 0)
	if bestCoversInt > 0 {
		for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers)); j++ {
			bestCoversForGenreResponse = append(bestCoversForGenreResponse, bestCovers[j])
		}
	}

	log.Debug().Msg("Genres got successfully")
	c.JSON(http.StatusOK, getResponse{
		GenreId:   genre.GenreId,
		Name:      genre.Name,
		BestCover: &bestCoversForGenreResponse,
	})
}
