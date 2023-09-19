package artist_service

import (
	"music-metadata/internal/clients/music_files_client/track_requests"
	"music-metadata/internal/database/repository"
)

type Service struct {
	ArtistRepo        repository.ArtistRepositoryInterface
	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     track_requests.TrackClient
}

func NewService(artistRepo repository.ArtistRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests track_requests.TrackClient) (service *Service) {

	service = &Service{
		ArtistRepo:        artistRepo,
		TrackMetadataRepo: trackMetadataRepo,
		TrackRequests:     trackRequests,
	}

	return service
}
