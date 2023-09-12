package album_service

import (
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type readAllResponseItem struct {
	AlbumId int    `json:"albumId"`
	Title   string `json:"title"`
}

type readAllResponse struct {
	Albums []readAllResponseItem `json:"albums"`
}

func (s *Service) ReadAll() (albums []models.Album, err error) {
	albums, err = s.AlbumRepo.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return make([]models.Album, 0), err
	}
	return albums, nil
}
