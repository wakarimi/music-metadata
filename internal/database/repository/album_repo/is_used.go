package album_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsUsed(tx *sqlx.Tx, albumId int) (used bool, err error) {
	query := `
		SELECT COUNT(*)
		FROM songs
		WHERE album_id = :album_id
	`

	args := map[string]interface{}{
		"album_id": albumId,
	}

	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to check if album is used")
		return false, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("Failed to scan count of songs for the album")
			return false, err
		}
	}

	return count > 0, nil
}
