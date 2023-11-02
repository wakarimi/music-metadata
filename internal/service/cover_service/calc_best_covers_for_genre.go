package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
)

func (s Service) CalcBestCoversForGenre(tx *sqlx.Tx, genreId int) (bestCovers []int, err error) {
	log.Debug().Int("genreId", genreId).Msg("Calculating best covers for genre")

	exists, err := s.SongService.GenreService.IsExists(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to check genre existence")
		return nil, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("genre with id=%d", genreId)}
		log.Error().Err(err).Int("genreId", genreId).Msg("Genre not found")
		return make([]int, 0), err
	}

	songs, err := s.SongService.GetAllByGenreId(tx, genreId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genre's songs")
		return make([]int, 0), err
	}

	audioFileIds := make([]int, len(songs))
	for i, song := range songs {
		audioFileIds[i] = song.AudioFileId
	}

	bestCovers, err = s.AudioFileClient.CoverTopForAudioFiles(audioFileIds)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch cover top")
		return make([]int, 0), err
	}

	log.Debug().Msg("Best covers calculated successfully successfully")
	return bestCovers, nil
}
