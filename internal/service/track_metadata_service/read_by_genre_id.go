package track_metadata_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TrackMetadataReadByGenreId struct {
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

func (s *Service) ReadByGenreId(tx *sqlx.Tx, genreId int) (trackMetadataList []TrackMetadataReadByGenreId, err error) {
	log.Debug().Int("genreId", genreId).Msg("Fetching track metadata list by genreId")

	trackMetadataModelList, err := s.TrackMetadataRepo.ReadAllByGenreTx(tx, genreId)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to fetch track metadata by genreId")
		return make([]TrackMetadataReadByGenreId, 0), err
	}

	trackMetadataList = make([]TrackMetadataReadByGenreId, len(trackMetadataModelList))
	for i, trackMetadataModel := range trackMetadataModelList {
		trackMetadataList[i] = TrackMetadataReadByGenreId{
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

	log.Debug().Int("genreId", genreId).Int("tracksCount", len(trackMetadataList)).Msg("Track metadata list by genreId fetched successfully")
	return trackMetadataList, nil
}
