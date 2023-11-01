package cover_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) CalcBestCoversForAlbum(tx *sqlx.Tx, albumId int) (bestCovers []int, err error) {
	log.Debug().Int("albumId", albumId).Msg("Calculating best covers for album")

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
