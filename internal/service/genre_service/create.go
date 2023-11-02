package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) Create(tx *sqlx.Tx, genre model.Genre) (createdGenre model.Genre, err error) {
	log.Debug().Interface("genre", genre).Msg("Creating new genre")

	genreId, err := s.GenreRepo.Create(tx, genre)
	if err != nil {
		log.Error().Err(err).Interface("genre", genre).Msg("Failed to create genre")
		return model.Genre{}, err
	}

	createdGenre, err = s.GenreRepo.Read(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to get created genre")
		return model.Genre{}, err
	}

	log.Debug().Interface("createdGenre", createdGenre).Msg("Genre created successfully")
	return createdGenre, nil
}
