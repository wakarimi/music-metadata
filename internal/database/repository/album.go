package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type AlbumRepositoryInterface interface {
	Create(album models.Album) (albumId int, err error)
	CreateTx(tx *sqlx.Tx, album models.Album) (albumId int, err error)
	Read(albumId int) (album models.Album, err error)
	ReadTx(tx *sqlx.Tx, albumId int) (album models.Album, err error)
	ReadByTitle(title string) (album models.Album, err error)
	ReadByTitleTx(tx *sqlx.Tx, title string) (album models.Album, err error)
	ReadAll() (albums []models.Album, err error)
	ReadAllTx(tx *sqlx.Tx) (albums []models.Album, err error)
	Delete(albumId int) (err error)
	DeleteTx(tx *sqlx.Tx, albumId int) (err error)
	IsExistsByTitle(title string) (exists bool, err error)
	IsExistsByTitleTx(tx *sqlx.Tx, title string) (exists bool, err error)
}

type AlbumRepository struct {
	Db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) AlbumRepositoryInterface {
	return &AlbumRepository{Db: db}
}

func (r *AlbumRepository) Create(album models.Album) (albumId int, err error) {
	log.Debug().Str("title", album.Title).Msg("Creating new album")
	return r.create(r.Db, album)
}

func (r *AlbumRepository) CreateTx(tx *sqlx.Tx, album models.Album) (albumId int, err error) {
	log.Debug().Str("title", album.Title).Msg("Creating new album transactional")
	return r.create(tx, album)
}

func (r *AlbumRepository) create(queryer Queryer, album models.Album) (albumId int, err error) {
	query := `
		INSERT INTO albums(title)
		VALUES (:title)
		RETURNING album_id
	`
	rows, err := queryer.NamedQuery(query, album)
	if err != nil {
		log.Error().Err(err).Str("title", album.Title).Msg("Failed to create album")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&albumId); err != nil {
			log.Error().Err(err).Str("title", album.Title).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after album insert")
		log.Error().Err(err).Str("title", album.Title).Msg("No id returned after album insert")
		return 0, err
	}

	log.Debug().Int("id", albumId).Str("title", album.Title).Msg("Album created successfully")
	return albumId, nil
}

func (r *AlbumRepository) Read(albumId int) (album models.Album, err error) {
	log.Debug().Int("id", albumId).Msg("Fetching album")
	return r.read(r.Db, albumId)
}

func (r *AlbumRepository) ReadTx(tx *sqlx.Tx, albumId int) (album models.Album, err error) {
	log.Debug().Int("id", albumId).Msg("Fetching album transactional")
	return r.read(tx, albumId)
}

func (r *AlbumRepository) read(queryer Queryer, albumId int) (album models.Album, err error) {
	query := `
		SELECT *
		FROM albums
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", albumId).Msg("Failed to fetch album")
		return models.Album{}, err
	}

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Err(err).Int("id", albumId).Msg("Failed to scan album into struct")
			return models.Album{}, err
		}
	} else {
		err := fmt.Errorf("no album found with id: %d", albumId)
		log.Error().Err(err).Int("id", albumId).Msg("No album found")
		return models.Album{}, err
	}

	log.Debug().Str("title", album.Title).Msg("Album fetched successfully")
	return album, nil
}

func (r *AlbumRepository) ReadByTitle(title string) (album models.Album, err error) {
	log.Debug().Str("title", title).Msg("Fetching album by title")
	return r.readByTitle(r.Db, title)
}

func (r *AlbumRepository) ReadByTitleTx(tx *sqlx.Tx, title string) (album models.Album, err error) {
	log.Debug().Str("title", title).Msg("Fetching album by title transactional")
	return r.readByTitle(tx, title)
}

func (r *AlbumRepository) readByTitle(queryer Queryer, title string) (album models.Album, err error) {
	query := `
		SELECT *
		FROM albums
		WHERE title = :title
	`
	args := map[string]interface{}{
		"title": title,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to fetch album")
		return models.Album{}, err
	}

	if rows.Next() {
		if err := rows.StructScan(&album); err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to scan album into struct")
			return models.Album{}, err
		}
	} else {
		err := fmt.Errorf("no album found with title: %s", title)
		log.Error().Err(err).Str("title", title).Msg("No album found")
		return models.Album{}, err
	}

	log.Debug().Int("id", album.AlbumId).Msg("Album fetched by title successfully")
	return album, nil
}

func (r *AlbumRepository) ReadAll() (albums []models.Album, err error) {
	log.Debug().Msg("Fetching all albums")
	return r.readAll(r.Db)
}

func (r *AlbumRepository) ReadAllTx(tx *sqlx.Tx) (albums []models.Album, err error) {
	log.Debug().Msg("Fetching all albums transactional")
	return r.readAll(tx)
}

func (r *AlbumRepository) readAll(queryer Queryer) (albums []models.Album, err error) {
	query := `
		SELECT *
		FROM albums
	`
	rows, err := queryer.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch albums")
		return make([]models.Album, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var album models.Album
		if err = rows.StructScan(&album); err != nil {
			log.Error().Err(err).Msg("Failed to scan albums data")
			return make([]models.Album, 0), err
		}
		albums = append(albums, album)
	}

	log.Debug().Int("count", len(albums)).Msg("All albums fetched successfully")
	return albums, nil
}

func (r *AlbumRepository) Delete(albumId int) (err error) {
	log.Debug().Int("id", albumId).Msg("Deleting album")
	return r.delete(r.Db, albumId)
}

func (r *AlbumRepository) DeleteTx(tx *sqlx.Tx, albumId int) (err error) {
	log.Debug().Int("id", albumId).Msg("Deleting album transactional")
	return r.delete(tx, albumId)
}

func (r *AlbumRepository) delete(queryer Queryer, albumId int) (err error) {
	query := `
		DELETE FROM albums
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", albumId).Msg("Failed to delete album")
		return err
	}

	log.Debug().Int("id", albumId).Msg("Album deleted successfully")
	return nil
}

func (r *AlbumRepository) IsExistsByTitle(title string) (exists bool, err error) {
	log.Debug().Str("title", title).Msg("Checking if album exists by title")
	return r.isExistsByTitle(r.Db, title)
}

func (r *AlbumRepository) IsExistsByTitleTx(tx *sqlx.Tx, title string) (exists bool, err error) {
	log.Debug().Str("title", title).Msg("Checking if album exists by title transactional")
	return r.isExistsByTitle(r.Db, title)
}

func (r *AlbumRepository) isExistsByTitle(queryer Queryer, title string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM albums
			WHERE title = :title
		)
	`
	args := map[string]interface{}{
		"title": title,
	}
	row, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to execute query to check album existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to scan result of album existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("title", title).Msg("Album exists")
	} else {
		log.Debug().Str("title", title).Msg("No album found")
	}
	return exists, nil
}
