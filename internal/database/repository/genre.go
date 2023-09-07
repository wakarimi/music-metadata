package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type GenreRepositoryInterface interface {
	CreateGenre(genre models.Genre) (genreId *int, err error)
	ReadGenre(genreId int) (genre models.Genre, err error)
	ReadGenreByName(name string) (genre models.Genre, err error)
	ReadAllGenres() (genres []models.Genre, err error)
	DeleteGenre(genreId int) error
	IsGenreExistsByName(name string) (bool, error)
}

type GenreRepository struct {
	Db *sqlx.DB
}

func NewGenreRepository(db *sqlx.DB) GenreRepositoryInterface {
	return &GenreRepository{Db: db}
}

func (r *GenreRepository) CreateGenre(genre models.Genre) (genreId *int, err error) {
	log.Info().Str("name", genre.Name).Msg("Creating new genre")

	query := `
		INSERT INTO genres(name)
		VALUES (:name)
		RETURNING genre_id
	`

	rows, err := r.Db.NamedQuery(query, genre)
	if err != nil {
		log.Error().Err(err).Str("name", genre.Name).Msg("Failed to create genre")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.Scan(&genreId); err != nil {
			log.Error().Err(err).Msg("Error scanning genreId from result set")
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no id returned after genre insert")
	}

	log.Info().Int("genreId", *genreId).Msg("Genre created successfully")
	return genreId, nil
}

func (r *GenreRepository) ReadGenre(genreId int) (genre models.Genre, err error) {
	log.Debug().Int("genreId", genreId).Msg("Fetching genre by ID")

	query := `
		SELECT *
		FROM genres
		WHERE genre_id = :genre_id
	`

	rows, err := r.Db.NamedQuery(query, map[string]interface{}{
		"genre_id": genreId,
	})
	if err != nil {
		log.Error().Int("genreId", genreId).Msg("Failed to fetch genre by ID")
		return models.Genre{}, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Int("genreId", genreId).Msg("Error scanning row into struct")
			return models.Genre{}, err
		}
	} else {
		return models.Genre{}, fmt.Errorf("no genre found with ID: %d", genreId)
	}

	log.Debug().Int("genreId", genreId).Msg("Fetched genre by ID successfully")
	return genre, nil
}

func (r *GenreRepository) ReadGenreByName(name string) (genre models.Genre, err error) {
	log.Debug().Str("name", name).Msg("Fetching genre by name")

	query := `
		SELECT *
		FROM genres
		WHERE name = :name
	`

	rows, err := r.Db.NamedQuery(query, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		log.Error().Str("name", name).Msg("Failed to fetch genre by name")
		return models.Genre{}, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Error closing rows")
		}
	}()

	if rows.Next() {
		if err := rows.StructScan(&genre); err != nil {
			log.Error().Str("name", name).Msg("Error scanning row into struct")
			return models.Genre{}, err
		}
	} else {
		return models.Genre{}, fmt.Errorf("no genre found with name: %s", name)
	}

	log.Debug().Str("name", name).Msg("Fetched genre by name successfully")
	return genre, nil
}

func (r *GenreRepository) ReadAllGenres() (genres []models.Genre, err error) {
	log.Info().Msg("Fetching all genres")

	query := `
		SELECT *
		FROM genres
	`
	err = r.Db.Select(&genres, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genres")
		return nil, err
	}

	log.Info().Int("count", len(genres)).Msg("Fetched all genres successfully")
	return genres, nil
}

func (r *GenreRepository) DeleteGenre(genreId int) error {
	log.Info().Int("genreId", genreId).Msg("Deleting genre")

	query := `
		DELETE FROM genres
		WHERE genre_id = :genre_id
	`

	result, err := r.Db.NamedExec(query, map[string]interface{}{
		"genre_id": genreId,
	})
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to delete genre")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to get rows affected after genre deletion")
		return err
	}
	if rowsAffected == 0 {
		log.Warn().Int("genreId", genreId).Msg("No rows affected while deleting genre")
	}

	log.Info().Int("genreId", genreId).Msg("Genre deleted successfully")
	return nil
}

func (r *GenreRepository) IsGenreExistsByName(name string) (bool, error) {
	log.Debug().Str("name", name).Msg("Checking if genre exists by name")

	var count int

	query := `
		SELECT COUNT(*)
		FROM genres
		WHERE name = :name
	`
	args := map[string]interface{}{
		"name": name,
	}

	rows, err := r.Db.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to check genre existence by name")
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
			log.Error().Err(err).Str("name", name).Msg("Failed to scan count from result set")
			return false, err
		}
	}

	return count > 0, nil
}
