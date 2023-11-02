package album_handler

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
	AlbumId   int    `json:"albumId"`
	Title     string `json:"title"`
	BestCover *[]int `json:"bestCovers,omitempty"`
}

func (h *Handler) Get(c *gin.Context) {
	log.Debug().Msg("Getting album")

	albumIdStr := c.Param("albumId")
	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		log.Error().Err(err).Str("albumIdStr", albumIdStr).Msg("Invalid dirId format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid albumId format",
			Reason:  err.Error(),
		})
		return
	}
	log.Debug().Int("albumId", albumId).Msg("Url parameter read successfully")

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

	var album model.Album
	var bestCovers []int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		album, err = h.AlbumService.Get(tx, albumId)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			bestCovers, err = h.CoverService.CalcBestCoversForAlbum(tx, album.AlbumId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get album")
		if _, ok := err.(errors.NotFound); ok {
			c.JSON(http.StatusNotFound, response.Error{
				Message: "Album not found",
				Reason:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to get album",
				Reason:  err.Error(),
			})
		}
		return
	}

	bestCoversForAlbumResponse := make([]int, 0)
	if bestCoversInt > 0 {
		for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers)); j++ {
			bestCoversForAlbumResponse = append(bestCoversForAlbumResponse, bestCovers[j])
		}
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, getResponse{
		AlbumId:   album.AlbumId,
		Title:     album.Title,
		BestCover: &bestCoversForAlbumResponse,
	})
}
