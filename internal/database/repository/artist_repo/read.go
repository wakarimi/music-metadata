package artist_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, artistId int) (artist model.Artist, err error) {
	query := `
		SELECT *
		FROM artists
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to fetch artist")
		return model.Artist{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Int("artistId", artistId).Msg("Failed to scan artist into struct")
			return model.Artist{}, err
		}
	} else {
		err := fmt.Errorf("no artist found with artist_id: %d", artistId)
		log.Error().Err(err).Int("artistId", artistId).Msg("No artist found")
		return model.Artist{}, err
	}

	log.Debug().Int("id", artist.ArtistId).Msg("Artist fetched successfully")
	return artist, nil
}
