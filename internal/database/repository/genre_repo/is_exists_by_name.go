package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByName(tx *sqlx.Tx, name string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM genres
			WHERE name = :name
		)
	`
	args := map[string]interface{}{
		"name": name,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to execute query to check genre existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan result of genre existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("name", name).Msg("Genre exists")
	} else {
		log.Debug().Str("name", name).Msg("No genre found")
	}
	return exists, nil
}
