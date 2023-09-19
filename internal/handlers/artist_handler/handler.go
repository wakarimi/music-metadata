package artist_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/cover_service"
	"music-metadata/internal/service/track_metadata_service"
)

type Handler struct {
	TransactionManager   service.TransactionManager
	ArtistService        artist_service.Service
	TrackMetadataService track_metadata_service.Service
	CoverService         cover_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	artistService artist_service.Service,
	trackMetadataService track_metadata_service.Service,
	coverService cover_service.Service) (h *Handler) {

	h = &Handler{
		TransactionManager:   transactionManager,
		ArtistService:        artistService,
		TrackMetadataService: trackMetadataService,
		CoverService:         coverService,
	}

	return h
}
