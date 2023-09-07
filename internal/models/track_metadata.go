package models

type TrackMetadata struct {
	TrackMetadataId int     `db:"track_metadata_id"`
	TrackId         int     `db:"track_id"`
	Title           *string `db:"title"`
	ArtistId        *int    `db:"artist_id"`
	AlbumId         *int    `db:"album_id"`
	Genre           *int    `db:"genre_id"`
	Bitrate         int     `db:"bitrate"`
	Channels        int     `db:"channels"`
	SampleRate      int     `db:"sample_rate"`
	Duration        int     `db:"duration"`
}
