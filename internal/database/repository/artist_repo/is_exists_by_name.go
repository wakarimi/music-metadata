package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByName(tx *sqlx.Tx, name string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM artists
			WHERE name = :name
		)
	`
	args := map[string]interface{}{
		"name": name,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to execute query to check artist existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan result of artist existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("name", name).Msg("Artist exists")
	} else {
		log.Debug().Str("name", name).Msg("No artist found")
	}
	return exists, nil
}
