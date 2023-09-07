package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"net/http"
)

type AlbumGetAllResponseOne struct {
	AlbumId int    `json:"albumId"`
	Title   string `json:"title"`
}

type AlbumGetAllResponse struct {
	Albums []AlbumGetAllResponseOne `json:"albums"`
}

func (h *AlbumHandler) GetAll(c *gin.Context) {
	log.Debug().Msg("Fetching all albums")

	albums, err := h.AlbumRepo.ReadAllAlbums()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all albums")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to get all albums",
		})
	}

	albumsResponse := make([]AlbumGetAllResponseOne, 0)
	for _, album := range albums {
		albumResponse := AlbumGetAllResponseOne{
			AlbumId: album.AlbumId,
			Title:   album.Title,
		}
		albumsResponse = append(albumsResponse, albumResponse)
	}

	log.Debug().Int("count", len(albumsResponse)).Msg("Fetched all albums successfully")
	c.JSON(http.StatusOK, AlbumGetAllResponse{
		Albums: albumsResponse,
	})
}
