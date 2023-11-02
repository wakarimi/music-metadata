package album_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) Get(tx *sqlx.Tx, albumId int) (album model.Album, err error) {
	log.Debug().Msg("Getting all albums")

	exists, err := s.IsExists(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to check existence")
		return model.Album{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("album with id=%d", albumId)}
		log.Error().Err(err).Int("albumId", albumId).Msg("Album not found")
		return model.Album{}, err
	}

	album, err = s.AlbumRepo.Read(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get albums")
		return model.Album{}, err
	}

	log.Debug().Interface("album", album).Msg("Albums got successfully")
	return album, nil
}
