package album_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
)

type Handler struct {
	AlbumService       album_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(albumService album_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		AlbumService:       albumService,
		TransactionManager: transactionManager,
	}

	return h
}
