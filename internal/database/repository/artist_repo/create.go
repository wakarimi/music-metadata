package artist_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Create(tx *sqlx.Tx, artist model.Artist) (artistId int, err error) {
	query := `
		INSERT INTO artists(name)
		VALUES (:name)
		RETURNING artist_id
	`
	rows, err := tx.NamedQuery(query, artist)
	if err != nil {
		log.Error().Err(err).Str("name", artist.Name).Msg("Failed to create artist")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&artistId); err != nil {
			log.Error().Err(err).Str("name", artist.Name).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after artist insert")
		log.Error().Err(err).Str("name", artist.Name).Msg("No id returned after artist insert")
		return 0, err
	}

	log.Debug().Int("id", artistId).Str("name", artist.Name).Msg("Artist created successfully")
	return artistId, nil
}
