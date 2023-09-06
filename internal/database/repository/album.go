package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"music-metadata/internal/models"
)

type AlbumRepositoryInterface interface {
	CreateAlbum(album models.Album) (albumId int, err error)
	ReadAlbum(albumId int) (album models.Album, err error)
	ReadAlbumByTitle(title string) (album models.Album, err error)
	ReadAllAlbums() ([]models.Album, error)
	DeleteAlbum(albumId int) error
	IsAlbumExistsByTitle(title string) (bool, error)
}

type AlbumRepository struct {
	Db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) AlbumRepositoryInterface {
	return &AlbumRepository{Db: db}
}

func (r *AlbumRepository) CreateAlbum(album models.Album) (albumId int, err error) {
	query := `
		INSERT INTO albums(title)
		VALUES (:title)
		RETURNING album_id
	`
	err = r.Db.QueryRow(query, album).Scan(&albumId)
	if err != nil {
		return 0, err
	}

	return albumId, nil
}

func (r *AlbumRepository) ReadAlbum(albumId int) (album models.Album, err error) {
	query := `
		SELECT album_id, title
		FROM albums
		WHERE album_id = :albumId
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Album{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&album, map[string]interface{}{
		"albumId": albumId,
	})
	if err != nil {
		return models.Album{}, err
	}

	return album, nil
}

func (r *AlbumRepository) ReadAlbumByTitle(title string) (album models.Album, err error) {
	query := `
		SELECT album_id, title
		FROM albums
		WHERE title = :title
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Album{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&album, map[string]interface{}{
		"title": title,
	})
	if err != nil {
		return models.Album{}, err
	}

	return album, nil
}

func (r *AlbumRepository) ReadAllAlbums() ([]models.Album, error) {
	query := `
		SELECT album_id, title
		FROM albums
	`

	var albums []models.Album
	err := r.Db.Select(&albums, query)
	if err != nil {
		log.Printf("Failed to fetch albums: %v", err)
		return nil, err
	}

	return albums, nil
}

func (r *AlbumRepository) DeleteAlbum(albumId int) error {
	query := `
		DELETE FROM albums
		WHERE album_id = :album_id
	`

	args := map[string]interface{}{
		"album_id": albumId,
	}

	_, err := r.Db.NamedExec(query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepository) IsAlbumExistsByTitle(title string) (bool, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM albums
		WHERE title = :title
	`

	args := map[string]interface{}{
		"title": title,
	}

	if err := r.Db.Get(&count, query, args); err != nil {
		return false, err
	}

	return count > 0, nil
}
