package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetAll(tx *sqlx.Tx) (albums []model.Album, err error) {
	log.Debug().Msg("Getting all albums")

	albums, err = s.AlbumRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get albums")
		return make([]model.Album, 0), err
	}

	log.Debug().Int("countOfAlbum", len(albums)).Msg("Albums got successfully")
	return albums, nil
}
