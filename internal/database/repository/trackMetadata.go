package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type TrackMetadataRepositoryInterface interface {
	CreateTrackMetadata(trackMetadata models.TrackMetadata) (trackMetadataId int, err error)
	ReadTrackMetadata(trackMetadataId int) (trackMetadata models.TrackMetadata, err error)
	ReadTrackMetadataByTrackId(trackId int) (trackMetadata models.TrackMetadata, err error)
	ReadAllTrackMetadata() (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTrackMetadataByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error)
	DeleteTrackMetadata(trackMetadataId int) error
	CountTracksByAlbum(albumId int) (count int, err error)
	CountTracksByArtist(artistId int) (count int, err error)
	CountTracksByGenre(genreId int) (count int, err error)
}

type TrackMetadataRepository struct {
	Db *sqlx.DB
}

func NewTrackMetadataRepository(db *sqlx.DB) TrackMetadataRepositoryInterface {
	return &TrackMetadataRepository{Db: db}
}

func (r *TrackMetadataRepository) CreateTrackMetadata(trackMetadata models.TrackMetadata) (trackMetadataId int, err error) {
	log.Info().Msg("Creating new track metadata")

	const query = `
		INSERT INTO track_metadata(track_id, title, artist_id, album_id, genre_id, bitrate, channels, sample_rate, duration)
		VALUES (:track_id :title :artist_id, :album_id, :genre_id, :bitrate, :channels, :sample_rate, :duration)
		RETURNING track_metadata_id
	`
	err = r.Db.QueryRow(query, trackMetadata).Scan(&trackMetadataId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create track metadata")
		return 0, err
	}

	log.Info().Int("trackMetadataId", trackMetadataId).Msg("Track metadata created successfully")
	return trackMetadataId, nil
}

func (r *TrackMetadataRepository) ReadTrackMetadata(trackMetadataId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("trackMetadataId", trackMetadataId).Msg("Fetching track metadata by ID")

	const query = `
		SELECT *
		FROM track_metadata
		WHERE track_metadata_id = :track_metadata_id
	`
	err = r.Db.Get(&trackMetadata, query, map[string]interface{}{
		"track_metadata_id": trackMetadataId,
	})
	if err != nil {
		log.Error().Err(err).Int("trackMetadataId", trackMetadataId).Msg("Failed to fetch track metadata by ID")
		return models.TrackMetadata{}, err
	}

	log.Debug().Int("trackMetadataId", trackMetadataId).Msg("Fetched track metadata by ID successfully")
	return trackMetadata, nil
}

func (r *TrackMetadataRepository) ReadTrackMetadataByTrackId(trackId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track metadata by track ID")

	const query = `
		SELECT *
		FROM track_metadata
		WHERE track_id = :track_id
	`
	err = r.Db.Get(&trackMetadata, query, map[string]interface{}{
		"track_id": trackId,
	})
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to fetch track metadata by track ID")
		return models.TrackMetadata{}, err
	}

	log.Debug().Int("trackId", trackId).Msg("Fetched track metadata by track ID successfully")
	return trackMetadata, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadata() (trackMetadataList []models.TrackMetadata, err error) {
	log.Info().Msg("Fetching all track metadata")

	const query = `
		SELECT *
		FROM track_metadata
	`
	err = r.Db.Select(&trackMetadataList, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return nil, err
	}

	log.Info().Int("count", len(trackMetadataList)).Msg("Fetched all track metadata successfully")
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching track metadata by album ID")

	const query = `
		SELECT *
		FROM track_metadata
		WHERE album_id = :album_id
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch track metadata by album ID")
		return nil, err
	}

	log.Debug().Int("count", len(trackMetadataList)).Int("albumId", albumId).Msg("Fetched track metadata by album ID successfully")
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("artistId", artistId).Msg("Fetching track metadata by artist ID")

	const query = `
		SELECT *
		FROM track_metadata
		WHERE artist_id = :artist_id
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"artist_id": artistId,
	})
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to fetch track metadata by artist ID")
		return nil, err
	}

	log.Debug().Int("count", len(trackMetadataList)).Int("artistId", artistId).Msg("Fetched track metadata by artist ID successfully")
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) ReadAllTrackMetadataByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("genreId", genreId).Msg("Fetching track metadata by genre ID")

	const query = `
		SELECT *
		FROM track_metadata
		WHERE genre_id = :genre_id
	`
	err = r.Db.Select(&trackMetadataList, query, map[string]interface{}{
		"genre_id": genreId,
	})
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to fetch track metadata by genre ID")
		return nil, err
	}

	log.Debug().Int("count", len(trackMetadataList)).Int("genreId", genreId).Msg("Fetched track metadata by genre ID successfully")
	return trackMetadataList, nil
}

func (r *TrackMetadataRepository) DeleteTrackMetadata(trackMetadataId int) error {
	log.Info().Int("trackMetadataId", trackMetadataId).Msg("Deleting track metadata")

	const query = `
		DELETE FROM track_metadata
		WHERE track_metadata_id = :track_metadata_id
	`
	_, err := r.Db.Exec(query, map[string]interface{}{
		"track_metadata_id": trackMetadataId,
	})
	if err != nil {
		log.Error().Err(err).Int("trackMetadataId", trackMetadataId).Msg("Failed to delete track metadata")
		return err
	}

	log.Info().Int("trackMetadataId", trackMetadataId).Msg("Track metadata deleted successfully")
	return nil
}

func (r *TrackMetadataRepository) CountTracksByAlbum(albumId int) (count int, err error) {
	log.Debug().Int("albumId", albumId).Msg("Counting tracks by album ID")

	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE album_id = :album_id
	`
	err = r.Db.Get(&count, query, map[string]interface{}{
		"album_id": albumId,
	})
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to count tracks by album ID")
		return 0, err
	}

	log.Debug().Int("count", count).Int("albumId", albumId).Msg("Counted tracks by album ID successfully")
	return count, nil
}

func (r *TrackMetadataRepository) CountTracksByArtist(artistId int) (count int, err error) {
	log.Debug().Int("artistId", artistId).Msg("Counting tracks by artist ID")

	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE artist_id = :artist_id
	`
	err = r.Db.Get(&count, query, map[string]interface{}{
		"artist_id": artistId,
	})
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to count tracks by artist ID")
		return 0, err
	}

	log.Debug().Int("count", count).Int("artistId", artistId).Msg("Counted tracks by artist ID successfully")
	return count, nil
}

func (r *TrackMetadataRepository) CountTracksByGenre(genreId int) (count int, err error) {
	log.Debug().Int("genreId", genreId).Msg("Counting tracks by genre ID")

	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE genre_id = :genre_id
	`
	err = r.Db.Get(&count, query, map[string]interface{}{
		"genre_id": genreId,
	})
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to count tracks by genre ID")
		return 0, err
	}

	log.Debug().Int("count", count).Int("genreId", genreId).Msg("Counted tracks by genre ID successfully")
	return count, nil
}
