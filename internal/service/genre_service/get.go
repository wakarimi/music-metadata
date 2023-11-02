package genre_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) Get(tx *sqlx.Tx, genreId int) (genre model.Genre, err error) {
	log.Debug().Msg("Getting all genres")

	exists, err := s.IsExists(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to check existence")
		return model.Genre{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("genre with id=%d", genreId)}
		log.Error().Err(err).Int("genreId", genreId).Msg("Genre not found")
		return model.Genre{}, err
	}

	genre, err = s.GenreRepo.Read(tx, genreId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genres")
		return model.Genre{}, err
	}

	log.Debug().Interface("genre", genre).Msg("Genres got successfully")
	return genre, nil
}
