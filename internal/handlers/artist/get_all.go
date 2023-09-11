package artist

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"net/http"
)

type ArtistGetAllResponseOne struct {
	ArtistId int    `json:"artistId"`
	Name     string `json:"name"`
}

type ArtistGetAllResponse struct {
	Artists []ArtistGetAllResponseOne `json:"artists"`
}

func (h *ArtistHandler) GetAll(c *gin.Context) {
	log.Debug().Msg("Fetching all artists")

	artists, err := h.ArtistRepo.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all artists")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all artists",
		})
	}

	artistsResponse := make([]ArtistGetAllResponseOne, 0)
	for _, artist := range artists {
		artistResponse := ArtistGetAllResponseOne{
			ArtistId: artist.ArtistId,
			Name:     artist.Name,
		}
		artistsResponse = append(artistsResponse, artistResponse)
	}

	log.Debug().Int("count", len(artistsResponse)).Msg("Fetched all artists successfully")
	c.JSON(http.StatusOK, ArtistGetAllResponse{
		Artists: artistsResponse,
	})
}
