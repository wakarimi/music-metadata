package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type ArtistRepositoryInterface interface {
	CreateArtist(artist models.Artist) (artistId int, err error)
	ReadArtist(artistId int) (artist models.Artist, err error)
	ReadArtistByName(name string) (artist models.Artist, err error)
	ReadAllArtists() (artists []models.Artist, err error)
	DeleteArtist(artistId int) error
	IsArtistExistsByName(name string) (bool, error)
}

type ArtistRepository struct {
	Db *sqlx.DB
}

func NewArtistRepository(db *sqlx.DB) ArtistRepositoryInterface {
	return &ArtistRepository{Db: db}
}

func (r *ArtistRepository) CreateArtist(artist models.Artist) (artistId int, err error) {
	log.Info().Str("name", artist.Name).Msg("Creating new artist")

	query := `
		INSERT INTO artists(name)
		VALUES (:name)
		RETURNING artist_id
	`
	err = r.Db.QueryRow(query, artist).Scan(&artistId)
	if err != nil {
		log.Error().Str("name", artist.Name).Msg("Failed to create artist")
		return 0, err
	}

	log.Info().Int("artistId", artistId).Msg("Artist created successfully")
	return artistId, nil
}

func (r *ArtistRepository) ReadArtist(artistId int) (artist models.Artist, err error) {
	log.Debug().Int("artistId", artistId).Msg("Fetching artist by ID")

	query := `
		SELECT *
		FROM artists
		WHERE artist_id = :artist_id
	`
	err = r.Db.Get(&artist, query, map[string]interface{}{
		"artist_id": artistId,
	})
	if err != nil {
		log.Error().Int("artistId", artistId).Msg("Failed to fetch artist by ID")
		return models.Artist{}, err
	}

	log.Debug().Int("artistId", artistId).Msg("Fetched artist by ID successfully")
	return artist, nil
}

func (r *ArtistRepository) ReadArtistByName(name string) (artist models.Artist, err error) {
	log.Debug().Str("name", name).Msg("Fetching artist by name")

	query := `
		SELECT *
		FROM artists
		WHERE name = :name
	`
	err = r.Db.Get(&artist, query, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		log.Error().Str("name", name).Msg("Failed to fetch artist by name")
		return models.Artist{}, err
	}

	log.Debug().Str("name", name).Msg("Fetched artist by name successfully")
	return artist, nil
}

func (r *ArtistRepository) ReadAllArtists() (artists []models.Artist, err error) {
	log.Info().Msg("Fetching all artists")

	query := `
		SELECT *
		FROM artists
	`
	err = r.Db.Select(&artists, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch artists")
		return nil, err
	}

	log.Info().Int("count", len(artists)).Msg("Fetched all artists successfully")
	return artists, nil
}

func (r *ArtistRepository) DeleteArtist(artistId int) error {
	log.Info().Int("artistId", artistId).Msg("Deleting artist")

	query := `
		DELETE FROM artists
		WHERE artist_id = :artist_id
	`
	_, err := r.Db.Exec(query, map[string]interface{}{
		"artist_id": artistId,
	})
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to delete artist")
		return err
	}

	log.Info().Int("artistId", artistId).Msg("Artist deleted successfully")
	return nil
}

func (r *ArtistRepository) IsArtistExistsByName(name string) (bool, error) {
	log.Debug().Str("name", name).Msg("Checking if artist exists by name")

	var count int

	query := `
		SELECT COUNT(*)
		FROM artists
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}
	if err := r.Db.Get(&count, query, args); err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to check artist existence by name")
		return false, err
	}

	return count > 0, nil
}
