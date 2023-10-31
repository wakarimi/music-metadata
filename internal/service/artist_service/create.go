package artist_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
)

func (s Service) Create(tx *sqlx.Tx, artist model.Artist) (createdArtist model.Artist, err error) {
	log.Debug().Interface("artist", artist).Msg("Creating new artist")

	artistId, err := s.ArtistRepo.Create(tx, artist)
	if err != nil {
		log.Error().Err(err).Interface("artist", artist).Msg("Failed to create artist")
		return model.Artist{}, err
	}

	createdArtist, err = s.ArtistRepo.Read(tx, artistId)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to get created artist")
		return model.Artist{}, err
	}

	log.Debug().Interface("createdArtist", createdArtist).Msg("Artist created successfully")
	return createdArtist, err
}
