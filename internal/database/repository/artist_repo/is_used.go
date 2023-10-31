package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsUsed(tx *sqlx.Tx, artistId int) (used bool, err error) {
	query := `
		SELECT COUNT(*)
		FROM songs
		WHERE artist_id = :artist_id
	`

	args := map[string]interface{}{
		"artist_id": artistId,
	}

	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to check if artist is used")
		return false, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("Failed to scan count of songs for the artist")
			return false, err
		}
	}

	return count > 0, nil
}
