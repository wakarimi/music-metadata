package genre_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) RemoveUnnecessaryItems(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Removing unnecessary items")

	genres, err := s.GenreRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read all genres")
		return err
	}

	for _, genre := range genres {
		used, err := s.GenreRepo.IsUsed(tx, genre.GenreId)
		if err != nil {
			log.Error().Err(err).Int("genreId", genre.GenreId).Msg("Failed to check usage")
			return err
		}

		if used {
			continue
		}

		err = s.GenreRepo.Delete(tx, genre.GenreId)
		if err != nil {
			log.Error().Err(err).Int("genreId", genre.GenreId).Msg("Failed to delete genre")
			return err
		}
	}

	log.Debug().Msg("Unnecessary items removed successfully")
	return nil
}
