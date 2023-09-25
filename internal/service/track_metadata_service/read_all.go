package track_metadata_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (genres []models.TrackMetadata, err error) {
	genres = make([]models.TrackMetadata, 0)

	genres, err = s.TrackMetadataRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return make([]models.TrackMetadata, 0), err
	}

	return genres, nil
}
