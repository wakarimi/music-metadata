package album_service

import (
	"music-metadata/internal/clients/musicfilesclient/trackrequests"
	"music-metadata/internal/database/repository"
)

type Service struct {
	AlbumRepo repository.AlbumRepositoryInterface

	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     trackrequests.TrackClient
}

func NewService(albumRepo repository.AlbumRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests trackrequests.TrackClient) (service *Service) {

	service = &Service{
		AlbumRepo:         albumRepo,
		TrackMetadataRepo: trackMetadataRepo,
		TrackRequests:     trackRequests,
	}

	return service
}
