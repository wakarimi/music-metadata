package genre_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Create(tx *sqlx.Tx, genre model.Genre) (genreId int, err error) {
	query := `
		INSERT INTO genres(name)
		VALUES (:name)
		RETURNING genre_id
	`
	rows, err := tx.NamedQuery(query, genre)
	if err != nil {
		log.Error().Err(err).Str("name", genre.Name).Msg("Failed to create genre")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&genreId); err != nil {
			log.Error().Err(err).Str("name", genre.Name).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after genre insert")
		log.Error().Err(err).Str("name", genre.Name).Msg("No id returned after genre insert")
		return 0, err
	}

	log.Debug().Int("id", genreId).Str("name", genre.Name).Msg("Genre created successfully")
	return genreId, nil
}
