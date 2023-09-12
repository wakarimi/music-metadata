package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/context"
	"music-metadata/internal/handlers/types"
	"net/http"
)

type readAllResponseItem struct {
	AlbumId     int    `json:"albumId"`
	Title       string `json:"title"`
	CoverId     *int   `json:"coverId"`
	TracksCount int    `json:"tracksCount"`
}

type readAllResponse struct {
	Albums []readAllResponseItem `json:"albums"`
}

func (h *Handler) ReadAll(ginCtx *gin.Context, appCtx *context.AppContext) {
	log.Debug().Msg("Fetching all albums")

	albums, err := h.AlbumService.ReadAll(appCtx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all albums")
		ginCtx.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch all albums",
		})
		return
	}

	albumsResponse := make([]readAllResponseItem, 0)
	for _, album := range albums {
		albumResponse := readAllResponseItem{
			AlbumId:     album.AlbumId,
			Title:       album.Title,
			CoverId:     album.CoverId,
			TracksCount: album.TracksCount,
		}
		albumsResponse = append(albumsResponse, albumResponse)
	}

	log.Debug().Int("count", len(albumsResponse)).Msg("All albums fetched successfully")
	ginCtx.JSON(http.StatusOK, readAllResponse{
		Albums: albumsResponse,
	})
}
