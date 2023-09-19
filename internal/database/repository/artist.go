package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type ArtistRepositoryInterface interface {
	Create(artist models.Artist) (artistId int, err error)
	CreateTx(tx *sqlx.Tx, artist models.Artist) (artistId int, err error)
	Read(artistId int) (artist models.Artist, err error)
	ReadTx(tx *sqlx.Tx, artistId int) (artist models.Artist, err error)
	ReadByName(name string) (artist models.Artist, err error)
	ReadByNameTx(tx *sqlx.Tx, name string) (artist models.Artist, err error)
	ReadAll() (artists []models.Artist, err error)
	ReadAllTx(tx *sqlx.Tx) (artists []models.Artist, err error)
	Delete(artistId int) (err error)
	DeleteTx(tx *sqlx.Tx, artistId int) (err error)
	IsExistsByName(name string) (exists bool, err error)
	IsExistsByNameTx(tx *sqlx.Tx, name string) (exists bool, err error)
}

type ArtistRepository struct {
	Db *sqlx.DB
}

func NewArtistRepository(db *sqlx.DB) ArtistRepositoryInterface {
	return &ArtistRepository{Db: db}
}

func (r *ArtistRepository) Create(artist models.Artist) (artistId int, err error) {
	log.Debug().Str("name", artist.Name).Msg("Creating new artist_handler")
	return r.create(r.Db, artist)
}

func (r *ArtistRepository) CreateTx(tx *sqlx.Tx, artist models.Artist) (artistId int, err error) {
	log.Debug().Str("name", artist.Name).Msg("Creating new artist_handler transactional")
	return r.create(tx, artist)
}

func (r *ArtistRepository) create(queryer Queryer, artist models.Artist) (artistId int, err error) {
	query := `
		INSERT INTO artists(name)
		VALUES (:name)
		RETURNING artist_id
	`
	rows, err := queryer.NamedQuery(query, artist)
	if err != nil {
		log.Error().Err(err).Str("name", artist.Name).Msg("Failed to create artist_handler")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&artistId); err != nil {
			log.Error().Err(err).Str("name", artist.Name).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after artist_handler insert")
		log.Error().Err(err).Str("name", artist.Name).Msg("No id returned after artist_handler insert")
		return 0, err
	}

	log.Debug().Int("id", artistId).Str("name", artist.Name).Msg("Artist created successfully")
	return artistId, nil
}

func (r *ArtistRepository) Read(artistId int) (artist models.Artist, err error) {
	log.Debug().Int("id", artistId).Msg("Fetching artist_handler")
	return r.read(r.Db, artistId)
}

func (r *ArtistRepository) ReadTx(tx *sqlx.Tx, artistId int) (artist models.Artist, err error) {
	log.Debug().Int("id", artistId).Msg("Fetching artist_handler transactional")
	return r.read(tx, artistId)
}

func (r *ArtistRepository) read(queryer Queryer, artistId int) (artist models.Artist, err error) {
	query := `
		SELECT *
		FROM artists
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", artistId).Msg("Failed to fetch artist_handler")
		return models.Artist{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Int("id", artistId).Msg("Failed to scan artist_handler into struct")
			return models.Artist{}, err
		}
	} else {
		err := fmt.Errorf("no artist_handler found with id: %d", artistId)
		log.Error().Err(err).Int("id", artistId).Msg("No artist_handler found")
		return models.Artist{}, err
	}

	log.Debug().Str("name", artist.Name).Msg("Artist fetched successfully")
	return artist, nil
}

func (r *ArtistRepository) ReadByName(name string) (artist models.Artist, err error) {
	log.Debug().Str("name", name).Msg("Fetching artist_handler by name")
	return r.readByName(r.Db, name)
}

func (r *ArtistRepository) ReadByNameTx(tx *sqlx.Tx, name string) (artist models.Artist, err error) {
	log.Debug().Str("name", name).Msg("Fetching artist_handler by name transactional")
	return r.readByName(tx, name)
}

func (r *ArtistRepository) readByName(queryer Queryer, name string) (artist models.Artist, err error) {
	query := `
		SELECT *
		FROM artists
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to fetch artist_handler")
		return models.Artist{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan artist_handler into struct")
			return models.Artist{}, err
		}
	} else {
		err := fmt.Errorf("no artist_handler found with name: %s", name)
		log.Error().Err(err).Str("name", name).Msg("No artist_handler found")
		return models.Artist{}, err
	}

	log.Debug().Int("id", artist.ArtistId).Msg("Artist fetched by name successfully")
	return artist, nil
}

func (r *ArtistRepository) ReadAll() (artists []models.Artist, err error) {
	log.Debug().Msg("Fetching all artists")
	return r.readAll(r.Db)
}

func (r *ArtistRepository) ReadAllTx(tx *sqlx.Tx) (artists []models.Artist, err error) {
	log.Debug().Msg("Fetching all artists transactional")
	return r.readAll(tx)
}

func (r *ArtistRepository) readAll(queryer Queryer) (artists []models.Artist, err error) {
	log.Debug().Msg("Fetching all artists")

	query := `
		SELECT *
		FROM artists
	`
	rows, err := queryer.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch artists")
		return make([]models.Artist, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var artist models.Artist
		if err = rows.StructScan(&artist); err != nil {
			log.Error().Err(err).Msg("Failed to scan artists data")
			return make([]models.Artist, 0), err
		}
		artists = append(artists, artist)
	}

	log.Debug().Int("count", len(artists)).Msg("All artists fetched successfully")
	return artists, nil
}

func (r *ArtistRepository) Delete(artistId int) (err error) {
	log.Debug().Int("id", artistId).Msg("Deleting artist_handler")
	return r.delete(r.Db, artistId)
}

func (r *ArtistRepository) DeleteTx(tx *sqlx.Tx, artistId int) (err error) {
	log.Debug().Int("id", artistId).Msg("Deleting artist_handler transactional")
	return r.delete(tx, artistId)
}

func (r *ArtistRepository) delete(queryer Queryer, artistId int) (err error) {
	log.Debug().Int("id", artistId).Msg("Deleting artist_handler")

	query := `
		DELETE FROM artists
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", artistId).Msg("Failed to delete artist_handler")
		return err
	}

	log.Debug().Int("id", artistId).Msg("Artist deleted successfully")
	return nil
}

func (r *ArtistRepository) IsExistsByName(name string) (exists bool, err error) {
	log.Debug().Str("name", name).Msg("Checking if artist_handler exists by name")
	return r.isExistsByName(r.Db, name)
}

func (r *ArtistRepository) IsExistsByNameTx(tx *sqlx.Tx, name string) (exists bool, err error) {
	log.Debug().Str("name", name).Msg("Checking if artist_handler exists by name transactional")
	return r.isExistsByName(r.Db, name)
}

func (r *ArtistRepository) isExistsByName(queryer Queryer, name string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM artists
			WHERE name = :name
		)
	`
	args := map[string]interface{}{
		"name": name,
	}
	row, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to execute query to check artist_handler existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan result of artist_handler existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("name", name).Msg("Artist exists")
	} else {
		log.Debug().Str("name", name).Msg("No artist_handler found")
	}
	return exists, nil
}
