package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type GenreRepositoryInterface interface {
	CreateGenre(genre models.Genre) (genreId int, err error)
	ReadGenre(genreId int) (genre models.Genre, err error)
	ReadGenreByName(name string) (genre models.Genre, err error)
	ReadAllGenres() ([]models.Genre, error)
	DeleteGenre(genreId int) error
	IsGenreExistsByName(name string) (bool, error)
}

type GenreRepository struct {
	Db *sqlx.DB
}

func NewGenreRepository(db *sqlx.DB) GenreRepositoryInterface {
	return &GenreRepository{Db: db}
}

func (r *GenreRepository) CreateGenre(genre models.Genre) (genreId int, err error) {
	query := `
		INSERT INTO genres(name)
		VALUES (:name)
		RETURNING genre_id
	`
	err = r.Db.QueryRow(query, genre).Scan(&genreId)
	if err != nil {
		return 0, err
	}

	return genreId, nil
}

func (r *GenreRepository) ReadGenre(genreId int) (genre models.Genre, err error) {
	query := `
		SELECT genre_id, name
		FROM genres
		WHERE genre_id = :genreId
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Genre{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&genre, map[string]interface{}{
		"genreId": genreId,
	})
	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (r *GenreRepository) ReadGenreByName(name string) (genre models.Genre, err error) {
	query := `
		SELECT genre_id, name
		FROM genres
		WHERE name = :name
	`

	namedStmt, err := r.Db.PrepareNamed(query)
	if err != nil {
		return models.Genre{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&genre, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (r *GenreRepository) ReadAllGenres() ([]models.Genre, error) {
	query := `
		SELECT genre_id, name
		FROM genres
	`

	var genres []models.Genre
	err := r.Db.Select(&genres, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch genres")
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) DeleteGenre(genreId int) error {
	query := `
		DELETE FROM genres
		WHERE genre_id = :genre_id
	`

	args := map[string]interface{}{
		"genre_id": genreId,
	}

	_, err := r.Db.NamedExec(query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *GenreRepository) IsGenreExistsByName(name string) (bool, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM genres
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
