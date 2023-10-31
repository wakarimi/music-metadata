package genre_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/genre_service"
)

type Handler struct {
	GenreService       genre_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(genreService genre_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		GenreService:       genreService,
		TransactionManager: transactionManager,
	}

	return h
}
