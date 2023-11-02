package genre_service

import "music-metadata/internal/database/repository/genre_repo"

type Service struct {
	GenreRepo genre_repo.Repo
}

func NewService(genreRepo genre_repo.Repo) (s *Service) {

	s = &Service{
		GenreRepo: genreRepo,
	}

	return s
}
