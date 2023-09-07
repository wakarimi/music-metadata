package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type AlbumRepositoryInterface interface {
	CreateAlbum(album models.Album) (albumId *int, err error)
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

func (r *AlbumRepository) CreateAlbum(album models.Album) (albumId *int, err error) {
	log.Info().Str("title", album.Title).Msg("Creating new album")

	query := `
		INSERT INTO albums(title)
		VALUES (:title)
		RETURNING album_id
	`

	rows, err := r.Db.NamedQuery(query, album)
	if err != nil {
		log.Error().Err(err).Str("title", album.Title).Msg("Failed to create album")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.Scan(&albumId); err != nil {
			log.Error().Err(err).Msg("Error scanning albumId from result set")
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no id returned after album insert")
	}

	log.Info().Int("albumId", *albumId).Msg("Album created successfully")
	return albumId, nil
}

func (r *AlbumRepository) ReadAlbum(albumId int) (album models.Album, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching album by ID")

	query := `
		SELECT *
		FROM albums
		WHERE album_id = :album_id
	`

	rows, err := r.Db.NamedQuery(query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Int("albumId", albumId).Msg("Failed to fetch album by ID")
		return models.Album{}, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Int("albumId", albumId).Msg("Error scanning row into struct")
			return models.Album{}, err
		}
	} else {
		return models.Album{}, fmt.Errorf("no album found with ID: %d", albumId)
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

	rows, err := r.Db.NamedQuery(query, map[string]interface{}{
		"title": title,
	})
	if err != nil {
		log.Error().Str("title", title).Msg("Failed to fetch album by title")
		return models.Album{}, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Str("title", title).Msg("Error scanning row into struct")
			return models.Album{}, err
		}
	} else {
		return models.Album{}, fmt.Errorf("no album found with title: %s", title)
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

	result, err := r.Db.NamedExec(query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to delete album")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to get rows affected after album deletion")
		return err
	}
	if rowsAffected == 0 {
		log.Warn().Int("albumId", albumId).Msg("No rows affected while deleting album")
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

	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to delete album")
		return false, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to scan count from result set")
			return false, err
		}
	}

	return count > 0, nil
}
