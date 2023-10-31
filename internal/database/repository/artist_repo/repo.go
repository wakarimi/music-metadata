package artist_repo

import (
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, artist model.Artist) (artistId int, err error)
	GetByName(tx *sqlx.Tx, name string) (artist model.Artist, err error)
	IsExistsByName(tx *sqlx.Tx, name string) (exists bool, err error)
	ReadAll(tx *sqlx.Tx) (artists []model.Artist, err error)
	Delete(tx *sqlx.Tx, artistId int) (err error)
	Read(tx *sqlx.Tx, artistId int) (artist model.Artist, err error)
	IsUsed(tx *sqlx.Tx, artistId int) (used bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
