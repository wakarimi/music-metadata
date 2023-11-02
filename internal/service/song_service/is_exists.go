package song_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, songId int) (exists bool, err error) {
	log.Debug().Int("songId", songId).Msg("Checking song existence")

	exists, err = s.SongRepo.IsExists(tx, songId)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to check song existence")
		return false, err
	}

	log.Debug().Int("songId", songId).Bool("exists", exists).Msg("Song existence checked successfully")
	return exists, nil
}
