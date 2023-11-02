package model

type Artist struct {
	ArtistId int    `db:"artist_id"`
	Name     string `db:"name"`
}
