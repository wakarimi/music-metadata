package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type AlbumRepositoryInterface interface {
	CreateAlbum(album models.Album) (albumId int, err error)
	ReadAlbum(albumId int) (album models.Album, err error)
	ReadAlbumByTitle(title string) (album models.Album, err error)
	ReadAllAlbums() (albums []models.Album, err error)
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
	log.Info().Str("title", album.Title).Msg("Creating new album")

	query := `
		INSERT INTO albums(title)
		VALUES (:title)
		RETURNING album_id
	`
	err = r.Db.QueryRow(query, album).Scan(&albumId)
	if err != nil {
		log.Error().Str("title", album.Title).Msg("Failed to create album")
		return 0, err
	}

	log.Info().Int("albumId", albumId).Msg("Album created successfully")
	return albumId, nil
}

func (r *AlbumRepository) ReadAlbum(albumId int) (album models.Album, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching album by ID")

	query := `
		SELECT *
		FROM albums
		WHERE album_id = :album_id
	`
	err = r.Db.Get(&album, query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Int("albumId", albumId).Msg("Failed to fetch album by ID")
		return models.Album{}, err
	}

	log.Debug().Int("albumId", albumId).Msg("Fetched album by ID successfully")
	return album, nil
}

func (r *AlbumRepository) ReadAlbumByTitle(title string) (album models.Album, err error) {
	log.Debug().Str("title", title).Msg("Fetching album by title")

	query := `
		SELECT *
		FROM albums
		WHERE title = :title
	`
	err = r.Db.Get(&album, query, map[string]interface{}{
		"title": title,
	})
	if err != nil {
		log.Error().Str("title", title).Msg("Failed to fetch album by title")
		return models.Album{}, err
	}

	log.Debug().Str("title", title).Msg("Fetched album by title successfully")
	return album, nil
}

func (r *AlbumRepository) ReadAllAlbums() (albums []models.Album, err error) {
	log.Info().Msg("Fetching all albums")

	query := `
		SELECT *
		FROM albums
	`
	err = r.Db.Select(&albums, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return nil, err
	}

	log.Info().Int("count", len(albums)).Msg("Fetched all albums successfully")
	return albums, nil
}

func (r *AlbumRepository) DeleteAlbum(albumId int) error {
	log.Info().Int("albumId", albumId).Msg("Deleting album")

	query := `
		DELETE FROM albums
		WHERE album_id = :album_id
	`
	_, err := r.Db.Exec(query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to delete album")
		return err
	}

	log.Info().Int("albumId", albumId).Msg("Album deleted successfully")
	return nil
}

func (r *AlbumRepository) IsAlbumExistsByTitle(title string) (bool, error) {
	log.Debug().Str("title", title).Msg("Checking if album exists by title")

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
		log.Error().Err(err).Str("title", title).Msg("Failed to check album existence by title")
		return false, err
	}

	return count > 0, nil
}
