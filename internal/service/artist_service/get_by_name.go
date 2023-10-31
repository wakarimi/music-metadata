package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetByName(tx *sqlx.Tx, name string) (artist model.Artist, err error) {
	log.Debug().Str("name", name).Msg("Getting new artist")

	artist, err = s.ArtistRepo.GetByName(tx, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to get artist")
		return model.Artist{}, err
	}

	log.Debug().Interface("artist", artist).Msg("Artist got successfully")
	return artist, err
}
