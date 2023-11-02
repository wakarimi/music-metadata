package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetByName(tx *sqlx.Tx, name string) (genre model.Genre, err error) {
	log.Debug().Str("name", name).Msg("Getting new genre")

	genre, err = s.GenreRepo.ReadByName(tx, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to get genre")
		return model.Genre{}, err
	}

	log.Debug().Interface("genre", genre).Msg("Genre got successfully")
	return genre, nil
}
