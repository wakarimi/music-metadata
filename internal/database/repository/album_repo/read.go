package album_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, albumId int) (album model.Album, err error) {
	query := `
		SELECT *
		FROM albums
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album")
		return model.Album{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Err(err).Int("albumId", albumId).Msg("Failed to scan album into struct")
			return model.Album{}, err
		}
	} else {
		err := fmt.Errorf("no album found with album_id: %d", albumId)
		log.Error().Err(err).Int("albumId", albumId).Msg("No album found")
		return model.Album{}, err
	}

	log.Debug().Int("id", album.AlbumId).Msg("Album fetched successfully")
	return album, nil
}
