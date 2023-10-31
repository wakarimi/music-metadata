package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"net/http"
)

func (h *Handler) Scan(c *gin.Context) {
	log.Debug().Msg("Scanning songs")

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		err = h.SongService.Scan(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan songs")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to scan songs",
			Reason:  err.Error(),
		})
		return
	}

	log.Debug().Msg("Songs scanned successfully")
	c.Status(http.StatusOK)
}
