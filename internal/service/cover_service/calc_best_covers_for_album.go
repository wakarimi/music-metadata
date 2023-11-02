package cover_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/errors"
)

func (s Service) CalcBestCoversForAlbum(tx *sqlx.Tx, albumId int) (bestCovers []int, err error) {
	log.Debug().Int("albumId", albumId).Msg("Calculating best covers for album")

	exists, err := s.SongService.AlbumService.IsExists(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to check album existence")
		return nil, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("album with id=%d", albumId)}
		log.Error().Err(err).Int("albumId", albumId).Msg("Album not found")
		return make([]int, 0), err
	}

	songs, err := s.SongService.GetAllByAlbumId(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get album's songs")
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
