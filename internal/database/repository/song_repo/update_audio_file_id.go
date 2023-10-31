package song_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) UpdateAudioFileId(tx *sqlx.Tx, songId int, audioFileId int) (err error) {
	query := `
        UPDATE songs
        SET audio_file_id = ?
        WHERE song_id = ?
    `

	result, err := tx.Exec(query, audioFileId, songId)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Int("audioFileId", audioFileId).Msg("Failed to update audio_file_id")
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

	log.Info().Int("songId", songId).Int("audioFileId", audioFileId).Msg("Audio file ID updated successfully")
	return nil
}
