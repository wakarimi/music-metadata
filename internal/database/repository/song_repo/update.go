package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Update(tx *sqlx.Tx, songId int, song model.Song) (err error) {
	query := `
		UPDATE songs
		SET audio_file_id = :audio_file_id, title = :title, album_id = :album_id, artist_id = :artist_id,
		    genre_id = :genre_id, year = :year, song_number = :song_number, disc_number = :disc_number,
		    lyrics = :lyrics, sha_256 = :sha_256
		WHERE song_id = :song_id
	`
	song.SongId = songId
	result, err := tx.NamedExec(query, song)
	if err != nil {
		log.Error().Err(err).Int("id", songId).Msg("Failed to update song")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Int("id", songId).Msg("Failed to get rows affected after song update")
		return err
	}
	if rowsAffected == 0 {
		err := fmt.Errorf("no rows affected while updating song")
		log.Error().Err(err).Int("id", songId).Msg("No rows affected while updating song")
		return err
	}

	log.Info().Int("id", songId).Msg("Song updated successfully")
	return nil
}
