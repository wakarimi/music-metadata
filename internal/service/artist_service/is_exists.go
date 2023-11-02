package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExists(tx *sqlx.Tx, artistId int) (exists bool, err error) {
	log.Debug().Int("artistId", artistId).Msg("Checking artist existence")

	exists, err = s.ArtistRepo.IsExists(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to check artist existence")
		return false, err
	}

	log.Debug().Int("artistId", artistId).Bool("exists", exists).Msg("Artist existence checked successfully")
	return exists, nil
}
