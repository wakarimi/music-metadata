package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) Read(tx *sqlx.Tx, albumId int) (genre models.Genre, err error) {
	genre, err = s.GenreRepo.ReadTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genre")
		return models.Genre{}, err
	}

	return genre, nil
}
