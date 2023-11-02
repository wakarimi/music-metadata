package album_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, albumId int) (err error) {
	query := `
		DELETE FROM albums
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", albumId).Msg("Failed to delete album")
		return err
	}

	log.Debug().Int("id", albumId).Msg("Album deleted successfully")
	return nil
}
