package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, genreId int) (exists bool, err error) {
	log.Debug().Int("genreId", genreId).Msg("Checking genre existence")

	exists, err = s.GenreRepo.IsExists(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to check genre existence")
		return false, err
	}

	log.Debug().Int("genreId", genreId).Bool("exists", exists).Msg("Genre existence checked successfully")
	return exists, nil
}
