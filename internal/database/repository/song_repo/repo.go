package song_repo

import (
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/model"
)

type Repo interface {
	ReadAll(tx *sqlx.Tx) (dirs []model.Song, err error)
	ReadAllByAlbumId(tx *sqlx.Tx, albumId int) (songs []model.Song, err error)
	ReadAllByArtistId(tx *sqlx.Tx, artistId int) (songs []model.Song, err error)
	ReadAllByGenreId(tx *sqlx.Tx, genreId int) (songs []model.Song, err error)
	Delete(tx *sqlx.Tx, songId int) (err error)
	Create(tx *sqlx.Tx, song model.Song) (songId int, err error)
	Update(tx *sqlx.Tx, songId int, song model.Song) (err error)
	UpdateAudioFileId(tx *sqlx.Tx, songId int, audioFileId int) (err error)
}

type Repository struct {
}

func NewRepository() Repo {
	return &Repository{}
}
