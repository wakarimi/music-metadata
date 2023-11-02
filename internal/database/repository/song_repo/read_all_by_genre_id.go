package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) ReadAllByGenreId(tx *sqlx.Tx, genreId int) (songs []model.Song, err error) {
	query := `
		SELECT *
		FROM songs
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch song")
		return make([]model.Song, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var song model.Song
		if err = rows.StructScan(&song); err != nil {
			log.Error().Err(err).Msg("Failed to scan song")
			return make([]model.Song, 0), err
		}
		songs = append(songs, song)
	}

	log.Debug().Int("genreId", genreId).Int("count", len(songs)).Msg("All song by genreId fetched successfully")
	return songs, nil
}
