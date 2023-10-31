package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) ReadAll(tx *sqlx.Tx) (artists []model.Artist, err error) {
	query := `
		SELECT *
		FROM artists
	`
	rows, err := tx.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch artists")
		return make([]model.Artist, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var artist model.Artist
		if err = rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Msg("Failed to scan artists data")
			return make([]model.Artist, 0), err
		}
		artists = append(artists, artist)
	}

	log.Debug().Int("count", len(artists)).Msg("All artists fetched successfully")
	return artists, nil
}
