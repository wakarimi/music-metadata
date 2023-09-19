package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) GetMostCommonCoverIdsByAlbumId(tx *sqlx.Tx, albumId int, n int) ([]int, error) {
	log.Debug().Int("albumId", albumId).Msg("Calculating most common coverIds for album")

	trackMetadataList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album's tracks")
		return nil, err
	}

	mostCommonCoverIds, err := s.getMostCommonCoverIds(trackMetadataList, n)
	if err != nil {
		log.Debug().Err(err).Int("albumId", albumId).Msg("Failed to get most common coverIds")
		return nil, err
	}

	if len(mostCommonCoverIds) > n {
		mostCommonCoverIds = mostCommonCoverIds[:n]
	}

	return mostCommonCoverIds, nil
}
