package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) Get(tx *sqlx.Tx, songId int) (song model.Song, err error) {
	log.Debug().Msg("Getting all songs")

	exists, err := s.IsExists(tx, songId)
	if err != nil {
		log.Error().Err(err).Int("songId", songId).Msg("Failed to check existence")
		return model.Song{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("song with id=%d", songId)}
		log.Error().Err(err).Int("songId", songId).Msg("Song not found")
		return model.Song{}, err
	}

	song, err = s.SongRepo.Read(tx, songId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get songs")
		return model.Song{}, err
	}

	log.Debug().Interface("song", song).Msg("Songs got successfully")
	return song, nil
}
