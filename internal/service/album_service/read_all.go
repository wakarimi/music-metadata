package album_service

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"io"
	"music-metadata/internal/context"
	"music-metadata/internal/models"
	"music-metadata/internal/service/music_files_responses"
	"net/http"
	"strconv"
)

type albumReadAll struct {
	AlbumId     int    `json:"albumId"`
	Title       string `json:"title"`
	CoverId     *int   `json:"coverId"`
	TracksCount int    `json:"tracksCount"`
}

func (s *Service) ReadAll(ctx *context.AppContext) (albums []albumReadAll, err error) {
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
				log.Error().Err(err).Msg("Failed to fetch album's tracks")
				return err
			}

			var coverId *int
			mostCommonCoverId, err := getMostCommonCoverId(ctx, trackMetadataList)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to get most common coverId")
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

func getMostCommonCoverId(ctx *context.AppContext, trackMetadataList []models.TrackMetadata) (mostCommonCoverId int, err error) {
	coverCounts := make(map[int]int)
	for _, trackMetadata := range trackMetadataList {
		resp, err := http.Get(ctx.Config.HttpServer.MusicFilesAddress + "/api/music-files-service/tracks/" + strconv.Itoa(trackMetadata.TrackId))
		if err != nil {
			log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Msg("Failed to fetch track")
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Str("status", resp.Status).Msg("Failed to fetch track")
			return 0, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Msg("Failed to read response body for track")
			return 0, err
		}

		var trackResponse music_files_responses.TrackGet
		err = json.Unmarshal(body, &trackResponse)
		if err != nil {
			log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Msg("Failed to deserialize response body for track")
			return 0, err
		}

		coverCounts[trackResponse.CoverId]++
	}

	mostCommonCoverId = 0
	maxCount := 0
	for coverId, count := range coverCounts {
		if count > maxCount {
			maxCount = count
			mostCommonCoverId = coverId
		}
	}
	if mostCommonCoverId == 0 {
		err := fmt.Errorf("covers not found")
		log.Debug().Err(err).Msg("Covers not found")
		return 0, err
	}
	return mostCommonCoverId, nil
}
