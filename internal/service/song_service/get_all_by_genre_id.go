package song_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) GetAllByGenreId(tx *sqlx.Tx, genreId int) (songs []model.Song, err error) {
	log.Debug().Int("genreId", genreId).Msg("Getting songs by genre")

	songs, err = s.SongRepo.ReadAllByGenreId(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to get songs by genre")
		return make([]model.Song, 0), err
	}

	log.Debug().Int("genreId", genreId).Int("countOfSongs", len(songs)).Msg("Songs by genre got successfully")
	return songs, nil
}
