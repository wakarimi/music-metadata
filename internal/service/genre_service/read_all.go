package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (genres []models.Genre, err error) {
	genres = make([]models.Genre, 0)

	genres, err = s.GenreRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genres")
		return make([]models.Genre, 0), err
	}

	return genres, nil
}
