package artist_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) Get(tx *sqlx.Tx, artistId int) (artist model.Artist, err error) {
	log.Debug().Msg("Getting all artists")

	exists, err := s.IsExists(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to check existence")
		return model.Artist{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("artist with id=%d", artistId)}
		log.Error().Err(err).Int("artistId", artistId).Msg("Artist not found")
		return model.Artist{}, err
	}

	artist, err = s.ArtistRepo.Read(tx, artistId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artists")
		return model.Artist{}, err
	}

	log.Debug().Interface("artist", artist).Msg("Artists got successfully")
	return artist, nil
}
