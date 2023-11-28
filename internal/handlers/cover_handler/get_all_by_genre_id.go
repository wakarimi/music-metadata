package cover_handler

import (
	"music-metadata/internal/handlers/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type getAllByGenreIdResponse struct {
	GenreID int   `json:"genreId"`
	Covers  []int `json:"covers"`
}

func (h *Handler) GetAllByGenreId(c *gin.Context) {
	log.Debug().Msg("Getting covers for genre")

	genreIDStr := c.Param("genreId")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		log.Error().Err(err).Str("genreIdStr", genreIDStr).Msg("Invalid genreId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid genreId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("genreId", genreID).Msg("Url parameter read successfully")

	var covers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		covers, err = h.CoverService.CalcBestCoversForGenre(tx, genreID)
		if err != nil {
			return err
		}
		return nil
	})

	log.Debug().Msg("Covers for genre got")
	c.JSON(http.StatusOK, getAllByGenreIdResponse{
		GenreID: genreID,
		Covers:  covers,
	})
}
