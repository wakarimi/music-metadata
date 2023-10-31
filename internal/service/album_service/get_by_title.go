package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetByTitle(tx *sqlx.Tx, title string) (album model.Album, err error) {
	log.Debug().Str("title", title).Msg("Getting new album")

	album, err = s.AlbumRepo.GetByTitle(tx, title)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to get album")
		return model.Album{}, err
	}

	log.Debug().Interface("album", album).Msg("Album got successfully")
	return album, err
}
