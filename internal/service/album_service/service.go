package album_service

import (
	"music-metadata/internal/database/repository"
	"music-metadata/internal/service"
)

type Service struct {
	TransactionManager service.TransactionManager
	AlbumRepo          repository.AlbumRepositoryInterface
	TrackMetadataRepo  repository.TrackMetadataRepositoryInterface
}

func NewService(transactionManager service.TransactionManager,
	albumRepo repository.AlbumRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface) *Service {

	return &Service{
		TransactionManager: transactionManager,
		AlbumRepo:          albumRepo,
		TrackMetadataRepo:  trackMetadataRepo,
	}
}
