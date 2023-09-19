package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) Read(tx *sqlx.Tx, albumId int) (album models.Album, err error) {
	album, err = s.AlbumRepo.ReadTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch album")
		return models.Album{}, err
	}

	return album, nil
}
