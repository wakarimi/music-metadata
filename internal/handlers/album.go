package handlers

import "music-metadata/internal/database/repository"

type AlbumHandler struct {
	AlbumRepo repository.AlbumRepositoryInterface
}

func NewAlbumHandler(
	albumRepo repository.AlbumRepositoryInterface) *AlbumHandler {
	return &AlbumHandler{albumRepo}
}
