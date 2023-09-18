package track_metadata_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TrackMetadataReadByAlbumId struct {
	TrackMetadataId int
	TrackId         int
	Title           *string
	ArtistId        *int
	AlbumId         *int
	GenreId         *int
	Year            *int
	TrackNumber     *int
	DiscNumber      *int
	Lyrics          *string
}

func (s *Service) ReadByAlbumId(tx *sqlx.Tx, albumId int) (trackMetadataList []TrackMetadataReadByAlbumId, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching track metadata list by albumId")

	trackMetadataModelList, err := s.TrackMetadataRepo.ReadAllByAlbumTx(tx, albumId)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch track metadata by albumId")
		return make([]TrackMetadataReadByAlbumId, 0), err
	}

	trackMetadataList = make([]TrackMetadataReadByAlbumId, len(trackMetadataModelList))
	for i, trackMetadataModel := range trackMetadataModelList {
		trackMetadataList[i] = TrackMetadataReadByAlbumId{
			TrackMetadataId: trackMetadataModel.TrackMetadataId,
			TrackId:         trackMetadataModel.TrackId,
			Title:           trackMetadataModel.Title,
			ArtistId:        trackMetadataModel.ArtistId,
			AlbumId:         trackMetadataModel.AlbumId,
			GenreId:         trackMetadataModel.GenreId,
			Year:            trackMetadataModel.Year,
			TrackNumber:     trackMetadataModel.TrackNumber,
			DiscNumber:      trackMetadataModel.DiscNumber,
			Lyrics:          trackMetadataModel.Lyrics,
		}
	}

	log.Debug().Int("albumId", albumId).Int("tracksCount", len(trackMetadataList)).Msg("Track metadata list by albumId fetched successfully")
	return trackMetadataList, nil
}
