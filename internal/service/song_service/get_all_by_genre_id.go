package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
	"music-metadata/internal/model"
)

func (s Service) GetAllByGenreId(tx *sqlx.Tx, genreId int) (songs []model.Song, err error) {
	log.Debug().Int("genreId", genreId).Msg("Getting songs by genre")

	exists, err := s.GenreService.IsExists(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to check genre existence")
		return nil, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("genre with id=%d", genreId)}
		log.Error().Err(err).Int("genreId", genreId).Msg("Genre not found")
		return make([]model.Song, 0), err
	}

	songs, err = s.SongRepo.ReadAllByGenreId(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to get songs by genre")
		return make([]model.Song, 0), err
	}

	log.Debug().Int("genreId", genreId).Int("countOfSongs", len(songs)).Msg("Songs by genre got successfully")
	return songs, nil
}
