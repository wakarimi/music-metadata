package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) RemoveUnnecessaryItems(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Removing unnecessary items")

	albums, err := s.AlbumRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read all albums")
		return err
	}

	for _, album := range albums {
		used, err := s.AlbumRepo.IsUsed(tx, album.AlbumId)
		if err != nil {
			log.Error().Err(err).Int("albumId", album.AlbumId).Msg("Failed to check usage")
			return err
		}

		if used {
			continue
		}

		err = s.AlbumRepo.Delete(tx, album.AlbumId)
		if err != nil {
			log.Error().Err(err).Int("albumId", album.AlbumId).Msg("Failed to delete album")
			return err
		}
	}

	log.Debug().Msg("Unnecessary items removed successfully")
	return nil
}
