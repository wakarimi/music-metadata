package model

type Genre struct {
	GenreId int    `db:"genre_id"`
	Name    string `db:"name"`
}
