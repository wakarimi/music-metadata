package artist_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) GetByName(tx *sqlx.Tx, name string) (artist model.Artist, err error) {
	query := `
		SELECT *
		FROM artists
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to fetch artist")
		return model.Artist{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan artist into struct")
			return model.Artist{}, err
		}
	} else {
		err := fmt.Errorf("no artist found with name: %s", name)
		log.Error().Err(err).Str("name", name).Msg("No artist found")
		return model.Artist{}, err
	}

	log.Debug().Int("id", artist.ArtistId).Msg("Artist fetched by name successfully")
	return artist, nil
}
