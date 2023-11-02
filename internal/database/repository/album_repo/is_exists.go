package album_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, albumId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM albums
			WHERE album_id = :album_id
		)
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to execute query to check album existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("albumId", albumId).Msg("Failed to scan result of album existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("albumId", albumId).Msg("Album exists")
	} else {
		log.Debug().Int("albumId", albumId).Msg("No album found")
	}
	return exists, nil
}
