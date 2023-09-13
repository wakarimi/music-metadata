package album_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/track_metadata_service"
)

type Handler struct {
	TransactionManager   service.TransactionManager
	AlbumService         album_service.Service
	TrackMetadataService track_metadata_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	albumService album_service.Service,
	trackMetadataService track_metadata_service.Service) *Handler {

	handler := &Handler{
		TransactionManager:   transactionManager,
		AlbumService:         albumService,
		TrackMetadataService: trackMetadataService,
	}

	return handler
}
