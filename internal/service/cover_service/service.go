package cover_service

import (
	"music-metadata/internal/clients/music_files_client/track_requests"
	"music-metadata/internal/database/repository"
)

type Service struct {
	AlbumRepo         repository.AlbumRepositoryInterface
	ArtistRepo        repository.ArtistRepositoryInterface
	GenreRepo         repository.GenreRepositoryInterface
	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     track_requests.TrackClient
}

func NewService(albumRepo repository.AlbumRepositoryInterface,
	artistRepo repository.ArtistRepositoryInterface,
	genreRepo repository.GenreRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests track_requests.TrackClient) (service *Service) {

	service = &Service{
		AlbumRepo:         albumRepo,
		ArtistRepo:        artistRepo,
		GenreRepo:         genreRepo,
		TrackMetadataRepo: trackMetadataRepo,
		TrackRequests:     trackRequests,
	}

	return service
}
