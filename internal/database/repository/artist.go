package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"music-metadata/internal/models"
)

type ArtistRepositoryInterface interface {
	CreateArtist(artist models.Artist) (artistId int, err error)
	ReadArtist(artistId int) (artist models.Artist, err error)
	ReadArtistByName(name string) (artist models.Artist, err error)
	ReadAllArtists() ([]models.Artist, error)
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
	query := `
		INSERT INTO artists(name)
		VALUES (:name)
		RETURNING artist_id
	`
	err = r.Db.QueryRow(query, artist).Scan(&artistId)
	if err != nil {
		return 0, err
	}

	return artistId, nil
}

func (r *ArtistRepository) ReadArtist(artistId int) (artist models.Artist, err error) {
	query := `
		SELECT artist_id, name
		FROM artists
		WHERE artist_id = :artistId
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Artist{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&artist, map[string]interface{}{
		"artistId": artistId,
	})
	if err != nil {
		return models.Artist{}, err
	}

	return artist, nil
}

func (r *ArtistRepository) ReadArtistByName(name string) (artist models.Artist, err error) {
	query := `
		SELECT artist_id, name
		FROM artists
		WHERE name = :name
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Artist{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&artist, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return models.Artist{}, err
	}

	return artist, nil
}

func (r *ArtistRepository) ReadAllArtists() ([]models.Artist, error) {
	query := `
		SELECT artist_id, name
		FROM artists
	`

	var artists []models.Artist
	err := r.Db.Select(&artists, query)
	if err != nil {
		log.Printf("Failed to fetch artists: %v", err)
		return nil, err
	}

	return artists, nil
}

func (r *ArtistRepository) DeleteArtist(artistId int) error {
	query := `
		DELETE FROM artists
		WHERE artist_id = :artist_id
	`

	args := map[string]interface{}{
		"artist_id": artistId,
	}

	_, err := r.Db.NamedExec(query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArtistRepository) IsArtistExistsByName(name string) (bool, error) {
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
		return false, err
	}

	return count > 0, nil
}
