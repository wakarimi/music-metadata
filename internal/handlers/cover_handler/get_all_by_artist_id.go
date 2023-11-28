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

type getAllByArtistIdResponse struct {
	ArtistID int   `json:"artistId"`
	Covers   []int `json:"covers"`
}

func (h *Handler) GetAllByArtistId(c *gin.Context) {
	log.Debug().Msg("Getting covers for artist")

	artistIDStr := c.Param("artistId")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Error().Err(err).Str("artistIdStr", artistIDStr).Msg("Invalid artistId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid artistId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("artistId", artistID).Msg("Url parameter read successfully")

	var covers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		covers, err = h.CoverService.CalcBestCoversForArtist(tx, artistID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artist's covers")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Artist not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get artist's covers",
				Reason:  err.Error(),
			})
		}
		return
	}

	log.Debug().Msg("Covers for artist got")
	c.JSON(http.StatusOK, getAllByArtistIdResponse{
		ArtistID: artistID,
		Covers:   covers,
	})
}
