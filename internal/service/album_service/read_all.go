package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type AlbumReadAll struct {
	AlbumId     int
	Title       string
	CoverId     *int
	TracksCount int
}

func (s *Service) ReadAll(tx *sqlx.Tx) (albums []AlbumReadAll, err error) {
	albums = make([]AlbumReadAll, 0)

	albumsModel, err := s.AlbumRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return make([]AlbumReadAll, 0), err
	}

	for _, albumModel := range albumsModel {
		trackMetadataList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumModel.AlbumId)
		if err != nil {
			log.Error().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to fetch album's tracks")
			return make([]AlbumReadAll, 0), err
		}

		var coverId *int
		mostCommonCoverId, err := s.getMostCommonCoverId(trackMetadataList)
		if err != nil {
			log.Debug().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to get most common coverId")
			coverId = nil
		} else {
			coverId = &mostCommonCoverId
		}

		album := AlbumReadAll{
			AlbumId:     albumModel.AlbumId,
			Title:       albumModel.Title,
			CoverId:     coverId,
			TracksCount: len(trackMetadataList),
		}
		albums = append(albums, album)
	}

	return albums, nil
}
