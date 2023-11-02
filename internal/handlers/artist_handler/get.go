package artist_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"modernc.org/mathutil"
	"music-metadata/internal/errors"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

type getResponse struct {
	ArtistId  int    `json:"artistId"`
	Name      string `json:"name"`
	BestCover *[]int `json:"bestCovers,omitempty"`
}

func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting artist")

	artistIdStr := c.Param("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil {
		log.Error().Err(err).Str("artistIdStr", artistIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid artistId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("artistId", artistId).Msg("Url parameter read successfully")

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

	var artist model.Artist
	var bestCovers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artist, err = h.ArtistService.Get(tx, artistId)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			bestCovers, err = h.CoverService.CalcBestCoversForArtist(tx, artist.ArtistId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artist")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Artist not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get artist",
				Reason:  err.Error(),
			})
		}
		return
	}

	bestCoversForArtistResponse := make([]int, 0)
	if bestCoversInt > 0 {
		for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers)); j++ {
			bestCoversForArtistResponse = append(bestCoversForArtistResponse, bestCovers[j])
		}
	}

	log.Debug().Msg("Artists got successfully")
	c.JSON(http.StatusOK, getResponse{
		ArtistId:  artist.ArtistId,
		Name:      artist.Name,
		BestCover: &bestCoversForArtistResponse,
	})
}
