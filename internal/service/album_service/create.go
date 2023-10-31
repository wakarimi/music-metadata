package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) Create(tx *sqlx.Tx, album model.Album) (createdAlbum model.Album, err error) {
	log.Debug().Interface("album", album).Msg("Creating new album")

	albumId, err := s.AlbumRepo.Create(tx, album)
	if err != nil {
		log.Error().Err(err).Interface("album", album).Msg("Failed to create album")
		return model.Album{}, err
	}

	createdAlbum, err = s.AlbumRepo.Read(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to get created album")
		return model.Album{}, err
	}

	log.Debug().Interface("createdAlbum", createdAlbum).Msg("Album created successfully")
	return createdAlbum, err
}
