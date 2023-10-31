package artist_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/artist_service"
)

type Handler struct {
	ArtistService      artist_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(artistService artist_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		ArtistService:      artistService,
		TransactionManager: transactionManager,
	}

	return h
}
