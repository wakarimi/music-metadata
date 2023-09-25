package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) Read(tx *sqlx.Tx, albumId int) (artist models.Artist, err error) {
	artist, err = s.ArtistRepo.ReadTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch artist")
		return models.Artist{}, err
	}

	return artist, nil
}
