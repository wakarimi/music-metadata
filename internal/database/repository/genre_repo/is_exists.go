package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, genreId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM genres
			WHERE genre_id = :genre_id
		)
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to execute query to check genre existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("genreId", genreId).Msg("Failed to scan result of genre existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("genreId", genreId).Msg("Genre exists")
	} else {
		log.Debug().Int("genreId", genreId).Msg("No genre found")
	}
	return exists, nil
}
