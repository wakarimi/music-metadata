package album_handler

import (
	"music-metadata/internal/service/album_service"
)

type Handler struct {
	AlbumService album_service.Service
}

func NewHandler(albumService album_service.Service) *Handler {
	return &Handler{
		albumService,
	}
}
