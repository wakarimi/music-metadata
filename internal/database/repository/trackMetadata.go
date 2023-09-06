package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"music-metadata/internal/models"
)

type TrackMetadataRepositoryInterface interface {
	CreateTrackMetadata(trackMetadata models.TrackMetadata) (trackMetadataId int, err error)
	ReadTrackMetadata(trackMetadataId int) (trackMetadata models.TrackMetadata, err error)
	ReadAllTrackMetadata() (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error)
	DeleteTrackMetadata(trackMetadataId int) error
}

type TrackMetadataRepository struct {
	Db *sqlx.DB
}

func NewTrackMetadataRepository(db *sqlx.DB) TrackMetadataRepositoryInterface {
	return &TrackMetadataRepository{Db: db}
}

func (r *TrackMetadataRepository) CreateTrackMetadata(trackMetadata models.TrackMetadata) (trackMetadataId int, err error) {
	query := `
		INSERT INTO track_metadata(artist_id, album_id, genre_id, bitrate, channels, sample_rate, duration)
		VALUES (:artist_id, :album_id, :genre_id, :bitrate, :channels, :sample_rate, :duration)
		RETURNING track_metadata_id
	`
	err = r.Db.QueryRow(query, trackMetadata).Scan(&trackMetadataId)
	if err != nil {
		return 0, err
	}
	return trackMetadataId, nil
}

func (r *TrackMetadataRepository) ReadTrackMetadata(trackMetadataId int) (trackMetadata models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE track_metadata_id = :trackMetadataId
	`
	err = r.Db.Get(&trackMetadata, query, map[string]interface{}{
		"trackMetadataId": trackMetadataId,
	})
	if err != nil {
		return models.TrackMetadata{}, err
	}
	return trackMetadata, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadata() (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
	`
	err = r.Db.Select(&trackMetadataList, query)
	if err != nil {
		log.Printf("Failed to fetch track metadata: %v", err)
		return nil, err
	}
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE album_id = :albumId
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"albumId": albumId,
	})
	if err != nil {
		return nil, err
	}
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE artist_id = :artistId
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"artistId": artistId,
	})
	if err != nil {
		return nil, err
	}
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE genre_id = :genreId
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"genreId": genreId,
	})
	if err != nil {
		return nil, err
	}
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) DeleteTrackMetadata(trackMetadataId int) error {
	query := `
		DELETE FROM track_metadata
		WHERE track_metadata_id = :trackMetadataId
	`
	_, err := r.Db.Exec(query, map[string]interface{}{
		"trackMetadataId": trackMetadataId,
	})
	if err != nil {
		return err
	}
	return nil
}
