package handlers

import (
	"music-metadata/internal/database/repository"
)

type MusicHandler struct {
	AlbumRepo  repository.AlbumRepositoryInterface
	ArtistRepo repository.ArtistRepositoryInterface
	GenreRepo  repository.GenreRepositoryInterface
	TrackRepo  repository.TrackMetadataRepositoryInterface
}

func NewMusicHandler(
	albumRepo repository.AlbumRepositoryInterface,
	artistRepo repository.ArtistRepositoryInterface,
	genreRepo repository.GenreRepositoryInterface,
	trackMetadataRepo repository.TrackMetadataRepositoryInterface) *MusicHandler {
	return &MusicHandler{albumRepo, artistRepo, genreRepo, trackMetadataRepo}
}
