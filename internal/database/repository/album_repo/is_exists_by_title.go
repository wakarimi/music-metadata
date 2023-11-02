package album_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByTitle(tx *sqlx.Tx, title string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM albums
			WHERE title = :title
		)
	`
	args := map[string]interface{}{
		"title": title,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to execute query to check album existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to scan result of album existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("title", title).Msg("Album exists")
	} else {
		log.Debug().Str("title", title).Msg("No album found")
	}
	return exists, nil
}
