package model

type Album struct {
	AlbumId int    `db:"album_id"`
	Title   string `db:"title"`
}
