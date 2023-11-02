package artist_service

import (
	"music-metadata/internal/database/repository/artist_repo"
)

type Service struct {
	ArtistRepo artist_repo.Repo
}

func NewService(artistRepo artist_repo.Repo) (s *Service) {

	s = &Service{
		ArtistRepo: artistRepo,
	}

	return s
}
