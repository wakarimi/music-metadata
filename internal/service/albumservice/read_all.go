package albumservice

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type albumReadAll struct {
	AlbumId     int
	Title       string
	CoverId     *int
	TracksCount int
}

func (s *Service) ReadAll() (albums []albumReadAll, err error) {
	albums = make([]albumReadAll, 0)

	err = s.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albumsModel, err := s.AlbumRepo.ReadAllTx(tx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch albums")
			return err
		}

		for _, albumModel := range albumsModel {
			trackMetadataList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumModel.AlbumId)
			if err != nil {
				log.Error().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to fetch album's tracks")
				return err
			}

			var coverId *int
			mostCommonCoverId, err := s.getMostCommonCoverId(trackMetadataList)
			if err != nil {
				log.Debug().Err(err).Int("albumId", albumModel.AlbumId).Msg("Failed to get most common coverId")
				coverId = nil
			} else {
				coverId = &mostCommonCoverId
			}

			album := albumReadAll{
				AlbumId:     albumModel.AlbumId,
				Title:       albumModel.Title,
				CoverId:     coverId,
				TracksCount: len(trackMetadataList),
			}
			albums = append(albums, album)
		}

		return nil
	})
	if err != nil {
		return make([]albumReadAll, 0), err
	}

	return albums, nil
}

func (s *Service) getMostCommonCoverId(trackMetadataList []models.TrackMetadata) (mostCommonCoverId int, err error) {
	log.Debug().Int("numberOfTracks", len(trackMetadataList)).Msg("Finding the most common cover")

	coverCounts, err := s.countCovers(trackMetadataList)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count covers counts")
		return 0, err
	}

	mostCommonCoverId, err = s.mostCommonCover(coverCounts)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to find most common cover")
		return 0, err
	}

	log.Debug().Int("mostCommonCoverId", mostCommonCoverId).Msg("Most common cover id found successfully")
	return mostCommonCoverId, nil
}

func (s *Service) countCovers(trackMetadataList []models.TrackMetadata) (coverCounts map[int]int, err error) {
	log.Debug().Msg("Counting cover encounters")

	coverCounts = make(map[int]int)

	for _, trackMetadata := range trackMetadataList {
		trackResponse, err := s.TrackRequests.GetById(trackMetadata.TrackId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch track")
			continue
		}

		if trackResponse.CoverId != nil {
			coverId := *trackResponse.CoverId
			coverCounts[coverId]++
		}
	}

	log.Debug().Int("numberOfDifferentCovers", len(coverCounts)).Msg("Covers counted successfully")
	return coverCounts, nil
}

func (s *Service) mostCommonCover(coverCounts map[int]int) (mostCommonCoverId int, err error) {
	log.Debug().Int("numberOfDifferentCovers", len(coverCounts)).Msg("Finding most common cover")

	if len(coverCounts) == 0 {
		err = fmt.Errorf("attempt to find the most common among 0 covers")
		log.Debug().Err(err).Msg("Attempt to find the most common among 0 covers")
		return 0, err
	}

	maxCount := 0

	for coverId, count := range coverCounts {
		if count > maxCount {
			maxCount = count
			mostCommonCoverId = coverId
		}
	}

	log.Debug().Int("coverId", mostCommonCoverId).Msg("Most common cover found successfully")
	return mostCommonCoverId, nil
}
