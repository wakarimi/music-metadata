package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (artists []models.Artist, err error) {
	artists = make([]models.Artist, 0)

	artists, err = s.ArtistRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch artists")
		return make([]models.Artist, 0), err
	}

	return artists, nil
}
