package genre

import "music-metadata/internal/database/repository"

type GenreHandler struct {
	GenreRepo repository.GenreRepositoryInterface
}

func NewGenreHandler(
	genreRepo repository.GenreRepositoryInterface) *GenreHandler {
	return &GenreHandler{genreRepo}
}
