package track_metadata_service

import (
	"music-metadata/internal/database/repository"
	"music-metadata/internal/service"
)

type Service struct {
	TransactionManager service.TransactionManager
	TrackMetadataRepo  repository.TrackMetadataRepositoryInterface
}

func NewService(transactionManager service.TransactionManager,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface) (service *Service) {

	service = &Service{
		TransactionManager: transactionManager,
		TrackMetadataRepo:  trackMetadataRepo,
	}

	return service
}
