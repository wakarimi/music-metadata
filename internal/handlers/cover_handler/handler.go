package cover_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/cover_service"
)

type Handler struct {
	CoverService       cover_service.Service
	TransactionManager service.TransactionManager
}

func NewHandler(coverService cover_service.Service,
	transactionManager service.TransactionManager,
) (h *Handler) {
	h = &Handler{
		CoverService:       coverService,
		TransactionManager: transactionManager,
	}

	return h
}
