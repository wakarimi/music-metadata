package album_service

import "music-metadata/internal/database/repository"

type Service struct {
	AlbumRepo repository.AlbumRepositoryInterface
}

func NewService(
	albumRepo repository.AlbumRepositoryInterface) *Service {
	return &Service{
		AlbumRepo: albumRepo,
	}
}
