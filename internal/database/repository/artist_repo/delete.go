package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, artistId int) (err error) {
	query := `
		DELETE FROM artists
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", artistId).Msg("Failed to delete artist")
		return err
	}

	log.Debug().Int("id", artistId).Msg("Artist deleted successfully")
	return nil
}
