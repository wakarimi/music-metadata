package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) GetMostCommonCoverIdsByGenreId(tx *sqlx.Tx, genreId int, n int) ([]int, error) {
	log.Debug().Int("genreId", genreId).Msg("Calculating most common coverIds for genre")

	trackMetadataList, err := s.TrackMetadataRepo.ReadAllByGenreTx(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to fetch genre's tracks")
		return nil, err
	}

	mostCommonCoverIds, err := s.getMostCommonCoverIds(trackMetadataList, n)
	if err != nil {
		log.Debug().Err(err).Int("genreId", genreId).Msg("Failed to get most common coverIds")
		return nil, err
	}

	if len(mostCommonCoverIds) > n {
		mostCommonCoverIds = mostCommonCoverIds[:n]
	}

	return mostCommonCoverIds, nil
}
