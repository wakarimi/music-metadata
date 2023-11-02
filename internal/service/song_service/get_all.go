package song_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetAll(tx *sqlx.Tx) (songs []model.Song, err error) {
	log.Debug().Msg("Getting all songs")

	songs, err = s.SongRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all songs")
		return make([]model.Song, 0), err
	}

	log.Debug().Int("countOfSongs", len(songs)).Msg("All songs got successfully")
	return songs, nil
}
