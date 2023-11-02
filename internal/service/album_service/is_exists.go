package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, albumId int) (exists bool, err error) {
	log.Debug().Int("albumId", albumId).Msg("Checking album existence")

	exists, err = s.AlbumRepo.IsExists(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to check album existence")
		return false, err
	}

	log.Debug().Int("albumId", albumId).Bool("exists", exists).Msg("Album existence checked successfully")
	return exists, nil
}
