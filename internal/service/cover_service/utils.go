package cover_service

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

func (s *Service) getMostCommonCoverIds(trackMetadataList []models.TrackMetadata, count int) (mostCommonCoverIds []int, err error) {
	log.Debug().Int("numberOfTracks", len(trackMetadataList)).Msg("Calculating the most common cover")

	coverCounts, err := s.countCovers(trackMetadataList)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count covers counts")
		return make([]int, 0), err
	}

	mostCommonCoverIds, err = s.mostCommonCovers(coverCounts, count)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to find most common cover")
		return make([]int, 0), err
	}

	log.Debug().Msg("Most common cover ids calculated successfully")
	return mostCommonCoverIds, nil
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

func (s *Service) mostCommonCovers(coverCounts map[int]int, n int) ([]int, error) {
	log.Debug().Int("numberOfDifferentCovers", len(coverCounts)).Msg("Finding most common covers")

	if len(coverCounts) == 0 {
		err := fmt.Errorf("attempt to find the most common among 0 covers")
		log.Debug().Err(err).Msg("Attempt to find the most common among 0 covers")
		return nil, err
	}

	var topCovers []int

	for i := 0; i < n && len(coverCounts) > 0; i++ {
		maxCover := getMaxCover(coverCounts)
		topCovers = append(topCovers, maxCover)
		delete(coverCounts, maxCover)
	}

	return topCovers, nil
}

func getMaxCover(coverCounts map[int]int) int {
	var maxCover int
	maxCount := -1

	for cover, count := range coverCounts {
		if count > maxCount {
			maxCover = cover
			maxCount = count
		}
	}

	return maxCover
}
