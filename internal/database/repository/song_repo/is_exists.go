package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, songId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM songs
			WHERE song_id = :song_id
		)
	`
	args := map[string]interface{}{
		"song_id": songId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to execute query to check song existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("songId", songId).Msg("Failed to scan result of song existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("songId", songId).Msg("Song exists")
	} else {
		log.Debug().Int("songId", songId).Msg("No song found")
	}
	return exists, nil
}
