package genre_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, genreId int) (genre model.Genre, err error) {
	query := `
		SELECT *
		FROM genres
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to fetch genre")
		return model.Genre{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Int("genreId", genreId).Msg("Failed to scan genre into struct")
			return model.Genre{}, err
		}
	} else {
		err := fmt.Errorf("no genre found with genre_id: %d", genreId)
		log.Error().Err(err).Int("genreId", genreId).Msg("No genre found")
		return model.Genre{}, err
	}

	log.Debug().Int("id", genre.GenreId).Msg("Genre fetched successfully")
	return genre, nil
}
