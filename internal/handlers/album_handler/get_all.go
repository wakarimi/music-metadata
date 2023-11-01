package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"modernc.org/mathutil"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

type getAllResponseItem struct {
	AlbumId   int    `json:"albumId"`
	Title     string `json:"title"`
	BestCover *[]int `json:"bestCovers,omitempty"`
}

type getAllResponse struct {
	Albums []getAllResponseItem `json:"albums"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting albums")

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

	var albums []model.Album
	var bestCovers [][]int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albums, err = h.AlbumService.GetAll(tx)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			for _, album := range albums {
				bestCoversItem, err := h.CoverService.CalcBestCoversForAlbum(tx, album.AlbumId)
				if err != nil {
					return err
				}
				bestCovers = append(bestCovers, bestCoversItem)
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get albums")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get albums",
			Reason:  err.Error(),
		})
		return
	}

	albumsResponseItems := make([]getAllResponseItem, len(albums))
	for i, album := range albums {
		bestCoversForAlbum := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers[i])); j++ {
				bestCoversForAlbum = append(bestCoversForAlbum, bestCovers[i][j])
			}
		}
		var bestCoverForAlbumResponse *[]int
		if len(bestCoversForAlbum) > 0 {
			bestCoverForAlbumResponse = &bestCoversForAlbum
		} else {
			bestCoverForAlbumResponse = nil
		}
		albumsResponseItems[i] = getAllResponseItem{
			AlbumId:   album.AlbumId,
			Title:     album.Title,
			BestCover: bestCoverForAlbumResponse,
		}
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Albums: albumsResponseItems,
	})
}
