package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type AlbumRead struct {
	AlbumId     int
	Title       string
	CoverId     *int
	TracksCount int
}

func (s *Service) Read(tx *sqlx.Tx, albumId int) (album AlbumRead, err error) {
	albumModel, err := s.AlbumRepo.ReadTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch album")
		return AlbumRead{}, err
	}

	trackMetadataList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumModel.AlbumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to fetch album's tracks")
		return AlbumRead{}, err
	}

	var coverId *int
	mostCommonCoverId, err := s.getMostCommonCoverId(trackMetadataList)
	if err != nil {
		log.Debug().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to get most common coverId")
		coverId = nil
	} else {
		coverId = &mostCommonCoverId
	}

	album = AlbumRead{
		AlbumId:     albumModel.AlbumId,
		Title:       albumModel.Title,
		CoverId:     coverId,
		TracksCount: len(trackMetadataList),
	}

	return album, nil
}
