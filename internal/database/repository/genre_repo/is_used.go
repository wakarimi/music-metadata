package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsUsed(tx *sqlx.Tx, genreId int) (used bool, err error) {
	query := `
		SELECT COUNT(*)
		FROM songs
		WHERE genre_id = :genre_id
	`

	args := map[string]interface{}{
		"genre_id": genreId,
	}

	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to check if genre is used")
		return false, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Error().Err(err).Msg("Failed to scan count of songs for the genre")
			return false, err
		}
	}

	return count > 0, nil
}
