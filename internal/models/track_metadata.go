package models

type TrackMetadata struct {
	TrackMetadataId int  `db:"track_metadata_id"`
	ArtistId        *int `db:"artist_id"`
	AlbumId         *int `db:"album_id"`
	Genre           *int `db:"genre_id"`
	Bitrate         int  `db:"bitrate"`
	Channels        int  `db:"channels"`
	SampleRate      int  `db:"sample_rate"`
	Duration        int  `db:"duration"`
}
