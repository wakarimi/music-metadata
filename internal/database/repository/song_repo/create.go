package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Create(tx *sqlx.Tx, song model.Song) (songId int, err error) {
	const query = `
		INSERT INTO songs(audio_file_id, title, album_id, artist_id, genre_id, year, song_number, disc_number, lyrics, sha_256)
		VALUES (:audio_file_id, :title, :album_id, :artist_id, :genre_id, :year, :song_number, :disc_number, :lyrics, :sha_256)
		RETURNING song_id
	`
	rows, err := tx.NamedQuery(query, song)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create song")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&songId); err != nil {
			log.Error().Err(err).Int("songId", song.SongId).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after song insert")
		log.Error().Err(err).Int("songId", song.SongId).Msg("No id returned after song insert")
		return 0, err
	}

	log.Info().Int("id", songId).Msg("Song created successfully")
	return songId, nil
}
