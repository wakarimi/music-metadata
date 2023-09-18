package album_service

import (
	"music-metadata/internal/clients/music_files_client/track_requests"
	"music-metadata/internal/database/repository"
)

type Service struct {
	AlbumRepo repository.AlbumRepositoryInterface

	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     track_requests.TrackClient
}

func NewService(albumRepo repository.AlbumRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests track_requests.TrackClient) (service *Service) {

	service = &Service{
		AlbumRepo:         albumRepo,
		TrackMetadataRepo: trackMetadataRepo,
		TrackRequests:     trackRequests,
	}

	return service
}
