package album_service

import (
	"music-metadata/internal/database/repository/album_repo"
)

type Service struct {
	AlbumRepo album_repo.Repo
}

func NewService(albumRepo album_repo.Repo) (s *Service) {

	s = &Service{
		AlbumRepo: albumRepo,
	}

	return s
}
