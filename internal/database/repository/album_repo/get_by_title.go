package album_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) GetByTitle(tx *sqlx.Tx, title string) (album model.Album, err error) {
	query := `
		SELECT *
		FROM albums
		WHERE title = :title
	`
	args := map[string]interface{}{
		"title": title,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to fetch album")
		return model.Album{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to scan album into struct")
			return model.Album{}, err
		}
	} else {
		err := fmt.Errorf("no album found with title: %s", title)
		log.Error().Err(err).Str("title", title).Msg("No album found")
		return model.Album{}, err
	}

	log.Debug().Int("id", album.AlbumId).Msg("Album fetched by title successfully")
	return album, nil
}
