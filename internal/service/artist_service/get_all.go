package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetAll(tx *sqlx.Tx) (artists []model.Artist, err error) {
	log.Debug().Msg("Getting all artists")

	artists, err = s.ArtistRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artists")
		return make([]model.Artist, 0), err
	}

	log.Debug().Int("countOfArtist", len(artists)).Msg("Artists got successfully")
	return artists, err
}
