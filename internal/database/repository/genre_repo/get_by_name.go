package genre_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) GetByName(tx *sqlx.Tx, name string) (genre model.Genre, err error) {
	query := `
		SELECT *
		FROM genres
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to fetch genre")
		return model.Genre{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan genre into struct")
			return model.Genre{}, err
		}
	} else {
		err := fmt.Errorf("no genre found with name: %s", name)
		log.Error().Err(err).Str("name", name).Msg("No genre found")
		return model.Genre{}, err
	}

	log.Debug().Int("id", genre.GenreId).Msg("Genre fetched by name successfully")
	return genre, nil
}
