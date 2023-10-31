package genre_repo

import (
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, genre model.Genre) (genreId int, err error)
	GetByName(tx *sqlx.Tx, name string) (genre model.Genre, err error)
	IsExistsByName(tx *sqlx.Tx, name string) (exists bool, err error)
	ReadAll(tx *sqlx.Tx) (genres []model.Genre, err error)
	Delete(tx *sqlx.Tx, genreId int) (err error)
	Read(tx *sqlx.Tx, genreId int) (genre model.Genre, err error)
	IsUsed(tx *sqlx.Tx, genreId int) (used bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}