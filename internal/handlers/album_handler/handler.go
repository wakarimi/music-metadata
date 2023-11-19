package album_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/cover_service"
)

type Handler struct {
	AlbumService       album_service.Service
	CoverService       cover_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(albumService album_service.Service,
	coverService cover_service.Service,
	transactionManager service.TransactionManager,
) (h *Handler) {
	h = &Handler{
		AlbumService:       albumService,
		CoverService:       coverService,
		TransactionManager: transactionManager,
	}

	return h
}
