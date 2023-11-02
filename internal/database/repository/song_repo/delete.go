package song_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, songId int) (err error) {
	query := `
		DELETE FROM songs
		WHERE song_id = :song_id
	`
	args := map[string]interface{}{
		"song_id": songId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to delete song")
		return err
	}

	log.Debug().Int("songId", songId).Msg("Song deleted successfully")
	return nil
}
