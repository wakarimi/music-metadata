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

type getAllResponseItem struct {
	GenreId   int    `json:"genreId"`
	Name      string `json:"name"`
	BestCover *[]int `json:"bestCovers,omitempty"`
}

type getAllResponse struct {
	Genres []getAllResponseItem `json:"genres"`
}

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
