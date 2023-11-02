package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, genreId int) (err error) {
	query := `
		DELETE FROM genres
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", genreId).Msg("Failed to delete genre")
		return err
	}

	log.Debug().Int("id", genreId).Msg("Genre deleted successfully")
	return nil
}
