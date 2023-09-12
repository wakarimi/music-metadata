package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/context"
	"music-metadata/internal/handlers/types"
	"net/http"
)

// readAllResponseItem godoc
// @Description Album details
// @Property AlbumId (integer) The unique identifier of the album
// @Property Title (string) Title of the album
// @Property CoverId (integer, optional) Identifier of the album cover
// @Property TracksCount (integer) Number of tracks in the album
type readAllResponseItem struct {
	AlbumId     int    `json:"albumId"`
	Title       string `json:"title"`
	CoverId     *int   `json:"coverId"`
	TracksCount int    `json:"tracksCount"`
}

// readAllResponse godoc
// @Description List of all albums
// @Property Albums (array) List of album details
type readAllResponse struct {
	Albums []readAllResponseItem `json:"albums"`
}

// ReadAll godoc
// @Summary Get all albums
// @Tags albums
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse
// @Failure 500 {object} types.Error
// @Router /albums [get]
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

	albumsResponse := make([]readAllResponseItem, len(albums))
	for i, album := range albums {
		albumResponse := readAllResponseItem{
			AlbumId:     album.AlbumId,
			Title:       album.Title,
			CoverId:     album.CoverId,
			TracksCount: album.TracksCount,
		}
		albumsResponse[i] = albumResponse
	}

	log.Debug().Int("count", len(albumsResponse)).Msg("All albums fetched successfully")
	ginCtx.JSON(http.StatusOK, readAllResponse{
		Albums: albumsResponse,
	})
}
