package models

type Song struct {
	SongId      int     `db:"song_id"`
	AudioFileId int     `db:"audio_file_id"`
	Title       *string `db:"title"`
	AlbumId     *int    `db:"album_id"`
	ArtistId    *int    `db:"artist_id"`
	GenreId     *int    `db:"genre_id"`
	Year        *int    `db:"year"`
	TrackNumber *int    `db:"track_number"`
	DiscNumber  *int    `db:"disc_number"`
	Lyrics      *string `db:"lyrics"`
	HashSha256  string  `db:"hash_sha_256"`
}
