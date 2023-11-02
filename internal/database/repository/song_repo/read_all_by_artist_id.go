package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) ReadAllByArtistId(tx *sqlx.Tx, artistId int) (songs []model.Song, err error) {
	query := `
		SELECT *
		FROM songs
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
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

	log.Debug().Int("artistId", artistId).Int("count", len(songs)).Msg("All song by artistId fetched successfully")
	return songs, nil
}
