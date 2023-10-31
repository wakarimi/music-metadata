package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) UpdateAudioFileId(tx *sqlx.Tx, songId int, audioFileId int) (err error) {
	query := `
		UPDATE songs
		SET audio_file_id = :audio_file_id
		WHERE song_id = :song_id
	`
	args := map[string]interface{}{
		"audio_file_id": audioFileId,
		"song_id":       songId,
	}
	result, err := tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to update song")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to get rows affected after audio_file_id update")
		return err
	}
	if rowsAffected == 0 {
		err := fmt.Errorf("no rows affected while updating audio_file_id")
		log.Error().Err(err).Int("songId", songId).Msg("No rows affected while updating audio_file_id")
		return err
	}

	log.Debug().Int("songId", songId).Msg("Song updated successfully")
	return nil
}
