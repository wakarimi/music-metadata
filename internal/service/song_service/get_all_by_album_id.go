package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) GetAllByAlbumId(tx *sqlx.Tx, albumId int) (songs []model.Song, err error) {
	log.Debug().Int("albumId", albumId).Msg("Getting songs by album")

	exists, err := s.AlbumService.IsExists(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to check album existence")
		return nil, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("album with id=%d", albumId)}
		log.Error().Err(err).Int("albumId", albumId).Msg("Album not found")
		return make([]model.Song, 0), err
	}

	songs, err = s.SongRepo.ReadAllByAlbumId(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to get songs by album")
		return make([]model.Song, 0), err
	}

	log.Debug().Int("albumId", albumId).Int("countOfSongs", len(songs)).Msg("Songs by album got successfully")
	return songs, nil
}
