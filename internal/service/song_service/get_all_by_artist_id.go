package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) GetAllByArtistId(tx *sqlx.Tx, artistId int) (songs []model.Song, err error) {
	log.Debug().Int("artistId", artistId).Msg("Getting songs by artist")

	exists, err := s.ArtistService.IsExists(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to check artist existence")
		return nil, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("artist with id=%d", artistId)}
		log.Error().Err(err).Int("artistId", artistId).Msg("Artist not found")
		return make([]model.Song, 0), err
	}

	songs, err = s.SongRepo.ReadAllByArtistId(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to get songs by artist")
		return make([]model.Song, 0), err
	}

	log.Debug().Int("artistId", artistId).Int("countOfSongs", len(songs)).Msg("Songs by artist got successfully")
	return songs, nil
}
