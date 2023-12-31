package cover_handler

import (
	"music-metadata/internal/errors"
	"music-metadata/internal/handlers/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type getAllByAlbumIdResponse struct {
	AlbumID int   `json:"albumId"`
	Covers  []int `json:"covers"`
}

func (h *Handler) GetAllByAlbumId(c *gin.Context) {
	log.Debug().Msg("Getting covers for album")

	albumIDStr := c.Param("albumId")
	albumID, err := strconv.Atoi(albumIDStr)
	if err != nil {
		log.Error().Err(err).Str("albumIdStr", albumIDStr).Msg("Invalid albumId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid albumId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("albumId", albumID).Msg("Url parameter read successfully")

	var covers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		covers, err = h.CoverService.CalcBestCoversForAlbum(tx, albumID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get album's covers")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Album not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get album's covers",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Covers for album got")
	c.JSON(http.StatusOK, getAllByAlbumIdResponse{
		AlbumID: albumID,
		Covers:  covers,
	})
}
