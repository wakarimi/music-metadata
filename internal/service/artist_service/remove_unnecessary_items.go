package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) RemoveUnnecessaryItems(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Removing unnecessary items")

	artists, err := s.ArtistRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read all artists")
		return err
	}

	for _, artist := range artists {
		used, err := s.ArtistRepo.IsUsed(tx, artist.ArtistId)
		if err != nil {
			log.Error().Err(err).Int("artistId", artist.ArtistId).Msg("Failed to check usage")
			return err
		}

		if used {
			continue
		}

		err = s.ArtistRepo.Delete(tx, artist.ArtistId)
		if err != nil {
			log.Error().Err(err).Int("artistId", artist.ArtistId).Msg("Failed to delete artist")
			return err
		}
	}

	log.Debug().Msg("Unnecessary items removed successfully")
	return nil
}
