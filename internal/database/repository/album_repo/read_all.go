package album_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (r Repository) ReadAll(tx *sqlx.Tx) (albums []model.Album, err error) {
	query := `
		SELECT *
		FROM albums
	`
	rows, err := tx.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return make([]model.Album, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var album model.Album
		if err = rows.StructScan(&album); err != nil {
			log.Error().Err(err).Msg("Failed to scan albums data")
			return make([]model.Album, 0), err
		}
		albums = append(albums, album)
	}

	log.Debug().Int("count", len(albums)).Msg("All albums fetched successfully")
	return albums, nil
}
