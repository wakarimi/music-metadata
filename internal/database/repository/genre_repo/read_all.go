package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) ReadAll(tx *sqlx.Tx) (genres []model.Genre, err error) {
	query := `
		SELECT *
		FROM genres
	`
	rows, err := tx.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genres")
		return make([]model.Genre, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var genre model.Genre
		if err = rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Msg("Failed to scan genres data")
			return make([]model.Genre, 0), err
		}
		genres = append(genres, genre)
	}

	log.Debug().Int("count", len(genres)).Msg("All genres fetched successfully")
	return genres, nil
}
