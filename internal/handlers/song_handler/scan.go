package song_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"net/http"
)

// Scan handles the request to initiate a scan for new or updated songs.
// @Summary Initiate a scan for new or updated songs
// @Description Scans the system for any new or updated songs, updating the database accordingly.
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 "Successfully initiated song scan"
// @Failure 500 {object} response.Error "Internal Server Error due to failure in scanning process"
// @Router /scan [post]
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
