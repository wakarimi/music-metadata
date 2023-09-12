package albumservice

import (
	"music-metadata/internal/clients/musicfilesclient/trackrequests"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/service"
)

type Service struct {
	TransactionManager service.TransactionManager
	AlbumRepo          repository.AlbumRepositoryInterface

	TrackMetadataRepo repository.TrackMetadataRepositoryInterface
	TrackRequests     trackrequests.TrackClient
}

func NewService(transactionManager service.TransactionManager,
	albumRepo repository.AlbumRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface,
	trackRequests trackrequests.TrackClient) (service *Service) {

	service = &Service{
		TransactionManager: transactionManager,
		AlbumRepo:          albumRepo,
		TrackMetadataRepo:  trackMetadataRepo,
		TrackRequests:      trackRequests,
	}

	return service
}
