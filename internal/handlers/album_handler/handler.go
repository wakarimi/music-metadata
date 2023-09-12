package album_handler

import (
	"music-metadata/internal/service/albumservice"
)

type Handler struct {
	AlbumService albumservice.Service
}

func NewHandler(albumService albumservice.Service) *Handler {
	return &Handler{
		albumService,
	}
}
