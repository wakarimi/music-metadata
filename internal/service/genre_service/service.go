package genre_service

import (
	"music-metadata/internal/clients/music_files_client/track_requests"
	"music-metadata/internal/database/repository"
)

type Service struct {
	GenreRepo         repository.GenreRepositoryInterface
	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     track_requests.TrackClient
}

func NewService(genreRepo repository.GenreRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests track_requests.TrackClient) (service *Service) {

	service = &Service{
		GenreRepo:         genreRepo,
		TrackMetadataRepo: trackMetadataRepo,
		TrackRequests:     trackRequests,
	}

	return service
}
