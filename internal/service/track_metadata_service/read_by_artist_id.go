package track_metadata_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TrackMetadataReadByArtistId struct {
	TrackMetadataId int
	TrackId         int
	Title           *string
	AlbumId         *int
	ArtistId        *int
	GenreId         *int
	Year            *int
	TrackNumber     *int
	DiscNumber      *int
	Lyrics          *string
}

func (s *Service) ReadByArtistId(tx *sqlx.Tx, artistId int) (trackMetadataList []TrackMetadataReadByArtistId, err error) {
	log.Debug().Int("artistId", artistId).Msg("Fetching track metadata list by artistId")

	trackMetadataModelList, err := s.TrackMetadataRepo.ReadAllByArtistTx(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to fetch track metadata by artistId")
		return make([]TrackMetadataReadByArtistId, 0), err
	}

	trackMetadataList = make([]TrackMetadataReadByArtistId, len(trackMetadataModelList))
	for i, trackMetadataModel := range trackMetadataModelList {
		trackMetadataList[i] = TrackMetadataReadByArtistId{
			TrackMetadataId: trackMetadataModel.TrackMetadataId,
			TrackId:         trackMetadataModel.TrackId,
			Title:           trackMetadataModel.Title,
			AlbumId:         trackMetadataModel.AlbumId,
			ArtistId:        trackMetadataModel.ArtistId,
			GenreId:         trackMetadataModel.GenreId,
			Year:            trackMetadataModel.Year,
			TrackNumber:     trackMetadataModel.TrackNumber,
			DiscNumber:      trackMetadataModel.DiscNumber,
			Lyrics:          trackMetadataModel.Lyrics,
		}
	}

	log.Debug().Int("artistId", artistId).Int("tracksCount", len(trackMetadataList)).Msg("Track metadata list by artistId fetched successfully")
	return trackMetadataList, nil
}
