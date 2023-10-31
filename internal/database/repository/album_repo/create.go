package album_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Create(tx *sqlx.Tx, album model.Album) (albumId int, err error) {
	query := `
		INSERT INTO albums(title)
		VALUES (:title)
		RETURNING album_id
	`
	rows, err := tx.NamedQuery(query, album)
	if err != nil {
		log.Error().Err(err).Str("title", album.Title).Msg("Failed to create album")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&albumId); err != nil {
			log.Error().Err(err).Str("title", album.Title).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after album insert")
		log.Error().Err(err).Str("title", album.Title).Msg("No id returned after album insert")
		return 0, err
	}

	log.Debug().Int("id", albumId).Str("title", album.Title).Msg("Album created successfully")
	return albumId, nil
}
