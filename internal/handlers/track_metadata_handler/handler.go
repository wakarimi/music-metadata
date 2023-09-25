package track_metadata_handler

import (
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/genre_service"
	"music-metadata/internal/service/track_metadata_service"
)

type Handler struct {
	TransactionManager   service.TransactionManager
	TrackMetadataService track_metadata_service.Service
	AlbumService         album_service.Service
	ArtistService        artist_service.Service
	GenreService         genre_service.Service
}

func NewHandler(transactionManager service.TransactionManager,
	trackMetadataService track_metadata_service.Service,
	albumService album_service.Service,
	artistService artist_service.Service,
	genreService genre_service.Service) (h *Handler) {

	h = &Handler{
		TransactionManager:   transactionManager,
		TrackMetadataService: trackMetadataService,
		AlbumService:         albumService,
		ArtistService:        artistService,
		GenreService:         genreService,
	}

	return h
}
