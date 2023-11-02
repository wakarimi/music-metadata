package album_repo

import (
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/model"
)

type Repo interface {
	Create(tx *sqlx.Tx, album model.Album) (albumId int, err error)
	Read(tx *sqlx.Tx, albumId int) (album model.Album, err error)
	ReadByTitle(tx *sqlx.Tx, title string) (album model.Album, err error)
	ReadAll(tx *sqlx.Tx) (albums []model.Album, err error)
	Delete(tx *sqlx.Tx, albumId int) (err error)
	IsExists(tx *sqlx.Tx, albumId int) (exists bool, err error)
	IsExistsByTitle(tx *sqlx.Tx, title string) (exists bool, err error)
	IsUsed(tx *sqlx.Tx, albumId int) (used bool, err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
