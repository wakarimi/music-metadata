package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExistsByName(tx *sqlx.Tx, name string) (exists bool, err error) {
	log.Debug().Str("name", name).Msg("Checking genre existence")

	exists, err = s.GenreRepo.IsExistsByName(tx, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to check genre existence")
		return false, err
	}

	log.Debug().Str("name", name).Bool("exists", exists).Msg("Genre existence checked successfully")
	return exists, err
}
