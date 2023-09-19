package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Service) GetMostCommonCoverIdsByArtistId(tx *sqlx.Tx, artistId int, n int) ([]int, error) {
	log.Debug().Int("artistId", artistId).Msg("Calculating most common coverIds for artist")

	trackMetadataList, err := s.TrackMetadataRepo.ReadAllByArtistTx(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to fetch artist's tracks")
		return nil, err
	}

	mostCommonCoverIds, err := s.getMostCommonCoverIds(trackMetadataList, n)
	if err != nil {
		log.Debug().Err(err).Int("artistId", artistId).Msg("Failed to get most common coverIds")
		return nil, err
	}

	if len(mostCommonCoverIds) > n {
		mostCommonCoverIds = mostCommonCoverIds[:n]
	}

	return mostCommonCoverIds, nil
}
