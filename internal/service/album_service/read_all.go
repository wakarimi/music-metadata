package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (albums []models.Album, err error) {
	albums = make([]models.Album, 0)

	albums, err = s.AlbumRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return make([]models.Album, 0), err
	}

	return albums, nil
}
