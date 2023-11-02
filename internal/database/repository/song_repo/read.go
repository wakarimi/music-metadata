package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, songId int) (song model.Song, err error) {
	query := `
		SELECT *
		FROM songs
		WHERE song_id = :song_id
	`
	args := map[string]interface{}{
		"song_id": songId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to fetch song")
		return model.Song{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&song); err != nil {
			log.Error().Err(err).Int("songId", songId).Msg("Failed to scan song into struct")
			return model.Song{}, err
		}
	} else {
		err := fmt.Errorf("no song found with song_id: %d", songId)
		log.Error().Err(err).Int("songId", songId).Msg("No song found")
		return model.Song{}, err
	}

	log.Debug().Int("id", song.SongId).Msg("Song fetched successfully")
	return song, nil
}
