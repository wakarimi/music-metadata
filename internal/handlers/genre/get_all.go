package genre

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"net/http"
)

type GenreGetAllResponseOne struct {
	GenreId int    `json:"genreId"`
	Name    string `json:"name"`
}

type GenreGetAllResponse struct {
	Genres []GenreGetAllResponseOne `json:"genres"`
}

func (h *GenreHandler) GetAll(c *gin.Context) {
	log.Debug().Msg("Fetching all genres")

	genres, err := h.GenreRepo.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all genres")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all genres",
		})
	}

	genresResponse := make([]GenreGetAllResponseOne, 0)
	for _, genre := range genres {
		genreResponse := GenreGetAllResponseOne{
			GenreId: genre.GenreId,
			Name:    genre.Name,
		}
		genresResponse = append(genresResponse, genreResponse)
	}

	log.Debug().Int("count", len(genresResponse)).Msg("Fetched all genres successfully")
	c.JSON(http.StatusOK, GenreGetAllResponse{
		Genres: genresResponse,
	})
}
