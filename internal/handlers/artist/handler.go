package artist

import "music-metadata/internal/database/repository"

type ArtistHandler struct {
	ArtistRepo repository.ArtistRepositoryInterface
}

func NewArtistHandler(
	artistRepo repository.ArtistRepositoryInterface) *ArtistHandler {
	return &ArtistHandler{artistRepo}
}
