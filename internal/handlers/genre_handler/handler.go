package genre_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/cover_service"
	"music-metadata/internal/service/genre_service"
)

type Handler struct {
	CoverService       cover_service.Service
	GenreService       genre_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(genreService genre_service.Service,
	coverService cover_service.Service,
	transactionManager service.TransactionManager) (h *Handler) {

	h = &Handler{
		GenreService:       genreService,
		CoverService:       coverService,
		TransactionManager: transactionManager,
	}

	return h
}
