package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, artistId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM artists
			WHERE artist_id = :artist_id
		)
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to execute query to check artist existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("artistId", artistId).Msg("Failed to scan result of artist existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("artistId", artistId).Msg("Artist exists")
	} else {
		log.Debug().Int("artistId", artistId).Msg("No artist found")
	}
	return exists, nil
}
