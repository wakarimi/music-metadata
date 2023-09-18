package models

type TrackMetadata struct {
	TrackMetadataId int     `db:"track_metadata_id"`
	TrackId         int     `db:"track_id"`
	Title           *string `db:"title"`
	AlbumId         *int    `db:"album_id"`
	ArtistId        *int    `db:"artist_id"`
	GenreId         *int    `db:"genre_id"`
	Year            *int    `db:"year"`
	TrackNumber     *int    `db:"track_number"`
	DiscNumber      *int    `db:"disc_number"`
	Lyrics          *string `db:"lyrics"`
	HashSha256      string  `db:"hash_sha_256"`
}
