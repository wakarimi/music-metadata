package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetAll(tx *sqlx.Tx) (genres []model.Genre, err error) {
	log.Debug().Msg("Getting all genres")

	genres, err = s.GenreRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genres")
		return make([]model.Genre, 0), err
	}

	log.Debug().Int("countOfGenre", len(genres)).Msg("Genres got successfully")
	return genres, err
}
