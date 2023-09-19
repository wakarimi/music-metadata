package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) GetMostCommonCoverIdByAlbumId(tx *sqlx.Tx, albumId int) (*int, error) {
	log.Debug().Int("albumId", albumId).Msg("Calculating most common coverId for album")

	trackMetadataList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album's tracks")
		return nil, err
	}

	mostCommonCoverId, err := s.getMostCommonCoverId(trackMetadataList)
	if err != nil {
		log.Debug().Err(err).Int("albumId", albumId).Msg("Failed to get most common coverId")
		return nil, err
	}

	log.Debug().Int("albumId", albumId).Int("mostCommonCoverId", mostCommonCoverId).Msg("Most common coverId for album calculated successfully")
	return &mostCommonCoverId, nil
}
