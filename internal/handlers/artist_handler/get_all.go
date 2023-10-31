package artist_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
)

type getAllResponseItem struct {
	ArtistId int    `json:"artistId"`
	Name     string `json:"name"`
}

type getAllResponse struct {
	Artists []getAllResponseItem `json:"artists"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting artists")

	var artists []model.Artist
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artists, err = h.ArtistService.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artists")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get artists",
			Reason:  err.Error(),
		})
		return
	}

	artistsResponseItems := make([]getAllResponseItem, len(artists))
	for i, artist := range artists {
		artistsResponseItems[i] = getAllResponseItem{
			ArtistId: artist.ArtistId,
			Name:     artist.Name,
		}
	}

	log.Debug().Msg("Artists got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Artists: artistsResponseItems,
	})
}
