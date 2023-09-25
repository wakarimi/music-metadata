package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type GenreRepositoryInterface interface {
	Create(genre models.Genre) (genreId int, err error)
	CreateTx(tx *sqlx.Tx, genre models.Genre) (genreId int, err error)
	Read(genreId int) (genre models.Genre, err error)
	ReadTx(tx *sqlx.Tx, genreId int) (genre models.Genre, err error)
	ReadByName(name string) (genre models.Genre, err error)
	ReadByNameTx(tx *sqlx.Tx, name string) (genre models.Genre, err error)
	ReadAll() (genres []models.Genre, err error)
	ReadAllTx(tx *sqlx.Tx) (genres []models.Genre, err error)
	Delete(genreId int) (err error)
	DeleteTx(tx *sqlx.Tx, genreId int) (err error)
	IsExistsByName(name string) (exists bool, err error)
	IsExistsByNameTx(tx *sqlx.Tx, name string) (exists bool, err error)
}

type GenreRepository struct {
	Db *sqlx.DB
}

func NewGenreRepository(db *sqlx.DB) GenreRepositoryInterface {
	return &GenreRepository{Db: db}
}

func (r *GenreRepository) Create(genre models.Genre) (genreId int, err error) {
	log.Debug().Str("name", genre.Name).Msg("Creating new genre")
	return r.create(r.Db, genre)
}

func (r *GenreRepository) CreateTx(tx *sqlx.Tx, genre models.Genre) (genreId int, err error) {
	log.Debug().Str("name", genre.Name).Msg("Creating new genre transactional")
	return r.create(tx, genre)
}

func (r *GenreRepository) create(queryer Queryer, genre models.Genre) (genreId int, err error) {
	query := `
		INSERT INTO genres(name)
		VALUES (:name)
		RETURNING genre_id
	`
	rows, err := queryer.NamedQuery(query, genre)
	if err != nil {
		log.Error().Err(err).Str("name", genre.Name).Msg("Failed to create genre")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&genreId); err != nil {
			log.Error().Err(err).Str("name", genre.Name).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after genre insert")
		log.Error().Err(err).Str("name", genre.Name).Msg("No id returned after genre insert")
		return 0, err
	}

	log.Debug().Int("id", genreId).Str("name", genre.Name).Msg("Genre created successfully")
	return genreId, nil
}

func (r *GenreRepository) Read(genreId int) (genre models.Genre, err error) {
	log.Debug().Int("id", genreId).Msg("Fetching genre")
	return r.read(r.Db, genreId)
}

func (r *GenreRepository) ReadTx(tx *sqlx.Tx, genreId int) (genre models.Genre, err error) {
	log.Debug().Int("id", genreId).Msg("Fetching genre transactional")
	return r.read(tx, genreId)
}

func (r *GenreRepository) read(queryer Queryer, genreId int) (genre models.Genre, err error) {
	query := `
		SELECT *
		FROM genres
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", genreId).Msg("Failed to fetch genre")
		return models.Genre{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Int("id", genreId).Msg("Failed to scan genre into struct")
			return models.Genre{}, err
		}
	} else {
		err := fmt.Errorf("no genre found with id: %d", genreId)
		log.Error().Err(err).Int("id", genreId).Msg("No genre found")
		return models.Genre{}, err
	}

	log.Debug().Str("name", genre.Name).Msg("Genre fetched successfully")
	return genre, nil
}

func (r *GenreRepository) ReadByName(name string) (genre models.Genre, err error) {
	log.Debug().Str("name", name).Msg("Fetching genre by name")
	return r.readByName(r.Db, name)
}

func (r *GenreRepository) ReadByNameTx(tx *sqlx.Tx, name string) (genre models.Genre, err error) {
	log.Debug().Str("name", name).Msg("Fetching genre by name transactional")
	return r.readByName(tx, name)
}

func (r *GenreRepository) readByName(queryer Queryer, name string) (genre models.Genre, err error) {
	query := `
		SELECT *
		FROM genres
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to fetch genre")
		return models.Genre{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan genre into struct")
			return models.Genre{}, err
		}
	} else {
		err := fmt.Errorf("no genre found with name: %s", name)
		log.Error().Err(err).Str("name", name).Msg("No genre found")
		return models.Genre{}, err
	}

	log.Debug().Int("id", genre.GenreId).Msg("Genre fetched by name successfully")
	return genre, nil
}

func (r *GenreRepository) ReadAll() (genres []models.Genre, err error) {
	log.Debug().Msg("Fetching all genres")
	return r.readAll(r.Db)
}

func (r *GenreRepository) ReadAllTx(tx *sqlx.Tx) (genres []models.Genre, err error) {
	log.Debug().Msg("Fetching all genres transactional")
	return r.readAll(tx)
}

func (r *GenreRepository) readAll(queryer Queryer) (genres []models.Genre, err error) {
	log.Debug().Msg("Fetching all genres")

	query := `
		SELECT *
		FROM genres
	`
	rows, err := queryer.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genres")
		return make([]models.Genre, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		if err = rows.StructScan(&genre); err != nil {
			log.Error().Err(err).Msg("Failed to scan genres data")
			return make([]models.Genre, 0), err
		}
		genres = append(genres, genre)
	}

	log.Debug().Int("count", len(genres)).Msg("All genres fetched successfully")
	return genres, nil
}

func (r *GenreRepository) Delete(genreId int) (err error) {
	log.Debug().Int("id", genreId).Msg("Deleting genre")
	return r.delete(r.Db, genreId)
}

func (r *GenreRepository) DeleteTx(tx *sqlx.Tx, genreId int) (err error) {
	log.Debug().Int("id", genreId).Msg("Deleting genre transactional")
	return r.delete(tx, genreId)
}

func (r *GenreRepository) delete(queryer Queryer, genreId int) (err error) {
	log.Debug().Int("id", genreId).Msg("Deleting genre")

	query := `
		DELETE FROM genres
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", genreId).Msg("Failed to delete genre")
		return err
	}

	log.Debug().Int("id", genreId).Msg("Genre deleted successfully")
	return nil
}

func (r *GenreRepository) IsExistsByName(name string) (exists bool, err error) {
	log.Debug().Str("name", name).Msg("Checking if genre exists by name")
	return r.isExistsByName(r.Db, name)
}

func (r *GenreRepository) IsExistsByNameTx(tx *sqlx.Tx, name string) (exists bool, err error) {
	log.Debug().Str("name", name).Msg("Checking if genre exists by name transactional")
	return r.isExistsByName(r.Db, name)
}

func (r *GenreRepository) isExistsByName(queryer Queryer, name string) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM genres
			WHERE name = :name
		)
	`
	args := map[string]interface{}{
		"name": name,
	}
	row, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to execute query to check genre existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to scan result of genre existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Str("name", name).Msg("Genre exists")
	} else {
		log.Debug().Str("name", name).Msg("No genre found")
	}
	return exists, nil
}
