package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
)

type getAllResponseItem struct {
	AlbumId int    `json:"albumId"`
	Title   string `json:"title"`
}

type getAllResponse struct {
	Albums []getAllResponseItem `json:"albums"`
}

func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting albums")

	var albums []model.Album
	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albums, err = h.AlbumService.GetAll(tx)
		if err != nil {
			return err
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
		albumsResponseItems[i] = getAllResponseItem{
			AlbumId: album.AlbumId,
			Title:   album.Title,
		}
	}

	log.Debug().Msg("Albums got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Albums: albumsResponseItems,
	})
}
