package genre_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
)

type getAllResponseItem struct {
	GenreId int    `json:"genreId"`
	Name    string `json:"name"`
}

type getAllResponse struct {
	Genres []getAllResponseItem `json:"genres"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting genres")

	var genres []model.Genre
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genres, err = h.GenreService.GetAll(tx)
		if err != nil {
			return err
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
		genresResponseItems[i] = getAllResponseItem{
			GenreId: genre.GenreId,
			Name:    genre.Name,
		}
	}

	log.Debug().Msg("Genres got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Genres: genresResponseItems,
	})
}
