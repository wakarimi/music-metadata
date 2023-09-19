package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/models"
)

type TrackMetadataRepositoryInterface interface {
	Create(trackMetadata models.TrackMetadata) (trackMetadataId int, err error)
	CreateTx(tx *sqlx.Tx, trackMetadata models.TrackMetadata) (trackMetadataId int, err error)
	Read(trackMetadataId int) (trackMetadata models.TrackMetadata, err error)
	ReadTx(tx *sqlx.Tx, trackMetadataId int) (trackMetadata models.TrackMetadata, err error)
	ReadByTrackId(trackId int) (trackMetadata models.TrackMetadata, err error)
	ReadByTrackIdTx(tx *sqlx.Tx, trackId int) (trackMetadata models.TrackMetadata, err error)
	ReadAll() (trackMetadataList []models.TrackMetadata, err error)
	ReadAllTx(tx *sqlx.Tx) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByAlbumTx(tx *sqlx.Tx, albumId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByArtistTx(tx *sqlx.Tx, artistId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error)
	ReadAllByGenreTx(tx *sqlx.Tx, genreId int) (trackMetadataList []models.TrackMetadata, err error)
	Update(trackMetadataId int, trackMetadata models.TrackMetadata) error
	UpdateTx(tx *sqlx.Tx, trackMetadataId int, trackMetadata models.TrackMetadata) error
	Delete(trackMetadataId int) (err error)
	DeleteTx(tx *sqlx.Tx, trackMetadataId int) (err error)
	CountByAlbum(albumId int) (count int, err error)
	CountByAlbumTx(tx *sqlx.Tx, albumId int) (count int, err error)
	CountByArtist(artistId int) (count int, err error)
	CountByArtistTx(tx *sqlx.Tx, artistId int) (count int, err error)
	CountByGenre(genreId int) (count int, err error)
	CountByGenreTx(tx *sqlx.Tx, genreId int) (count int, err error)
	IsExistsByTrackId(trackId int) (exists bool, err error)
	IsExistsByTrackIdTx(tx *sqlx.Tx, trackId int) (exists bool, err error)
}

type TrackMetadataRepository struct {
	Db *sqlx.DB
}

func NewTrackMetadataRepository(db *sqlx.DB) TrackMetadataRepositoryInterface {
	return &TrackMetadataRepository{Db: db}
}

func (r TrackMetadataRepository) Create(trackMetadata models.TrackMetadata) (trackMetadataId int, err error) {
	log.Debug().Int("trackId", trackMetadata.TrackId).Msg("Creating new track metadata")
	return r.create(r.Db, trackMetadata)
}

func (r TrackMetadataRepository) CreateTx(tx *sqlx.Tx, trackMetadata models.TrackMetadata) (trackMetadataId int, err error) {
	log.Debug().Int("trackId", trackMetadata.TrackId).Msg("Creating new track metadata transactional")
	return r.create(tx, trackMetadata)
}

func (r TrackMetadataRepository) create(queryer Queryer, trackMetadata models.TrackMetadata) (trackMetadataId int, err error) {
	const query = `
		INSERT INTO track_metadata(track_id, title, album_id, artist_id, genre_id, year, track_number, disc_number, lyrics, hash_sha_256)
		VALUES (:track_id, :title, :album_id, :artist_id, :genre_id, :year, :track_number, :disc_number, :lyrics, :hash_sha_256)
		RETURNING track_metadata_id
	`
	rows, err := queryer.NamedQuery(query, trackMetadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create track metadata")
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&trackMetadataId); err != nil {
			log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Msg("Failed to scan id into filed")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after track metadata insert")
		log.Error().Err(err).Int("trackId", trackMetadata.TrackId).Msg("No id returned after track metadata insert")
		return 0, err
	}

	log.Info().Int("id", trackMetadataId).Msg("Track metadata created successfully")
	return trackMetadataId, nil
}

func (r TrackMetadataRepository) Read(trackMetadataId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("id", trackMetadataId).Msg("Fetching track metadata by id")
	return r.read(r.Db, trackMetadataId)
}

func (r TrackMetadataRepository) ReadTx(tx *sqlx.Tx, trackMetadataId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("id", trackMetadataId).Msg("Fetching track metadata by id transactional")
	return r.read(tx, trackMetadataId)
}

func (r TrackMetadataRepository) read(queryer Queryer, trackMetadataId int) (trackMetadata models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE track_metadata_id = :track_metadata_id
	`
	args := map[string]interface{}{
		"track_metadata_id": trackMetadataId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", trackMetadataId).Msg("Failed to fetch track metadata")
		return models.TrackMetadata{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Int("id", trackMetadataId).Msg("Failed to scan track metadata into struct")
			return models.TrackMetadata{}, err
		}
	} else {
		err := fmt.Errorf("no track metadata found with id: %d", trackMetadataId)
		log.Error().Err(err).Int("id", trackMetadataId).Msg("No track metadata found")
		return models.TrackMetadata{}, err
	}

	log.Debug().Int("trackId", trackMetadata.TrackId).Msg("Track metadata fetched successfully")
	return trackMetadata, nil
}

func (r TrackMetadataRepository) ReadByTrackId(trackId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track metadata by trackId")
	return r.readByTrackId(r.Db, trackId)
}

func (r TrackMetadataRepository) ReadByTrackIdTx(tx *sqlx.Tx, trackId int) (trackMetadata models.TrackMetadata, err error) {
	log.Debug().Int("trackId", trackId).Msg("Fetching track metadata by trackId transactional")
	return r.readByTrackId(tx, trackId)
}

func (r TrackMetadataRepository) readByTrackId(queryer Queryer, trackId int) (trackMetadata models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE track_id = :track_id
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to fetch track metadata")
		return models.TrackMetadata{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Int("trackId", trackId).Msg("Failed to scan track metadata into struct")
			return models.TrackMetadata{}, err
		}
	} else {
		err := fmt.Errorf("no track metadata found with track_id: %d", trackId)
		log.Error().Err(err).Int("trackId", trackId).Msg("No track metadata found")
		return models.TrackMetadata{}, err
	}

	log.Debug().Int("id", trackMetadata.TrackMetadataId).Msg("Track metadata fetched successfully")
	return trackMetadata, nil
}

func (r TrackMetadataRepository) ReadAll() (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Msg("Fetching all track metadata")
	return r.readAll(r.Db)
}

func (r TrackMetadataRepository) ReadAllTx(tx *sqlx.Tx) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Msg("Fetching all track metadata transactional")
	return r.readAll(tx)
}

func (r TrackMetadataRepository) readAll(queryer Queryer) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
	`
	rows, err := queryer.Queryx(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return make([]models.TrackMetadata, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var trackMetadata models.TrackMetadata
		if err = rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Msg("Failed to scan track metadata data")
			return make([]models.TrackMetadata, 0), err
		}
		trackMetadataList = append(trackMetadataList, trackMetadata)
	}

	log.Debug().Int("count", len(trackMetadataList)).Msg("All track metadata fetched successfully")
	return trackMetadataList, nil
}

func (r TrackMetadataRepository) ReadAllByAlbum(albumId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching all track metadata by albumId")
	return r.readAllByAlbum(r.Db, albumId)
}

func (r TrackMetadataRepository) ReadAllByAlbumTx(tx *sqlx.Tx, albumId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("albumId", albumId).Msg("Fetching all track metadata by albumId transactional")
	return r.readAllByAlbum(tx, albumId)
}

func (r TrackMetadataRepository) readAllByAlbum(queryer Queryer, albumId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return make([]models.TrackMetadata, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var trackMetadata models.TrackMetadata
		if err = rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Msg("Failed to scan track metadata data")
			return make([]models.TrackMetadata, 0), err
		}
		trackMetadataList = append(trackMetadataList, trackMetadata)
	}

	log.Debug().Int("albumId", albumId).Int("count", len(trackMetadataList)).Msg("All track metadata by albumId fetched successfully")
	return trackMetadataList, nil
}

func (r TrackMetadataRepository) ReadAllByArtist(artistId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("artistId", artistId).Msg("Fetching all track metadata by artistId")
	return r.readAllByArtist(r.Db, artistId)
}

func (r TrackMetadataRepository) ReadAllByArtistTx(tx *sqlx.Tx, artistId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("artistId", artistId).Msg("Fetching all track metadata by artistId transactional")
	return r.readAllByArtist(tx, artistId)
}

func (r TrackMetadataRepository) readAllByArtist(queryer Queryer, artistId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return make([]models.TrackMetadata, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var trackMetadata models.TrackMetadata
		if err = rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Msg("Failed to scan track metadata data")
			return make([]models.TrackMetadata, 0), err
		}
		trackMetadataList = append(trackMetadataList, trackMetadata)
	}

	log.Debug().Int("artistId", artistId).Int("count", len(trackMetadataList)).Msg("All track metadata by artistId fetched successfully")
	return trackMetadataList, nil
}

func (r TrackMetadataRepository) ReadAllByGenre(genreId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("genreId", genreId).Msg("Fetching all track metadata by artistId")
	return r.readAllByGenre(r.Db, genreId)
}

func (r TrackMetadataRepository) ReadAllByGenreTx(tx *sqlx.Tx, genreId int) (trackMetadataList []models.TrackMetadata, err error) {
	log.Debug().Int("genreId", genreId).Msg("Fetching all track metadata by genreId transactional")
	return r.readAllByGenre(tx, genreId)
}

func (r TrackMetadataRepository) readAllByGenre(queryer Queryer, genreId int) (trackMetadataList []models.TrackMetadata, err error) {
	query := `
		SELECT *
		FROM track_metadata
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	rows, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch track metadata")
		return make([]models.TrackMetadata, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var trackMetadata models.TrackMetadata
		if err = rows.StructScan(&trackMetadata); err != nil {
			log.Error().Err(err).Msg("Failed to scan track metadata data")
			return make([]models.TrackMetadata, 0), err
		}
		trackMetadataList = append(trackMetadataList, trackMetadata)
	}

	log.Debug().Int("genreId", genreId).Int("count", len(trackMetadataList)).Msg("All track metadata by genre fetched successfully")
	return trackMetadataList, nil
}

func (r TrackMetadataRepository) Update(trackMetadataId int, trackMetadata models.TrackMetadata) error {
	log.Debug().Int("id", trackMetadataId).Msg("Updating track metadata")
	return r.update(r.Db, trackMetadataId, trackMetadata)
}

func (r TrackMetadataRepository) UpdateTx(tx *sqlx.Tx, trackMetadataId int, trackMetadata models.TrackMetadata) error {
	log.Debug().Int("id", trackMetadataId).Msg("Updating track metadata transactional")
	return r.update(tx, trackMetadataId, trackMetadata)
}

func (r TrackMetadataRepository) update(queryer Queryer, trackMetadataId int, trackMetadata models.TrackMetadata) error {
	const query = `
		UPDATE track_metadata
		SET track_id = :track_id, title = :title, album_id = :album_id, artist_id = :artist_id, genre_id = :genre_id, 
		    year = :year, track_number = :track_number, disc_number = :disc_number, lyrics = :lyrics,
		    hash_sha_256 = :hash_sha_256
		WHERE track_metadata_id = :track_metadata_id
	`
	trackMetadata.TrackMetadataId = trackMetadataId
	result, err := queryer.NamedExec(query, trackMetadata)
	if err != nil {
		log.Error().Err(err).Int("id", trackMetadataId).Msg("Failed to update track metadata")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Int("id", trackMetadataId).Msg("Failed to get rows affected after track metadata update")
		return err
	}
	if rowsAffected == 0 {
		err := fmt.Errorf("no rows affected while updating track metadata")
		log.Error().Err(err).Int("id", trackMetadataId).Msg("No rows affected while updating track metadata")
		return err
	}

	log.Info().Int("id", trackMetadataId).Msg("Track metadata updated successfully")
	return nil
}

func (r TrackMetadataRepository) Delete(trackMetadataId int) (err error) {
	log.Debug().Int("id", trackMetadataId).Msg("Deleting track metadata")
	return r.delete(r.Db, trackMetadataId)
}

func (r TrackMetadataRepository) DeleteTx(tx *sqlx.Tx, trackMetadataId int) (err error) {
	log.Debug().Int("id", trackMetadataId).Msg("Deleting track metadata transactional")
	return r.delete(tx, trackMetadataId)
}

func (r TrackMetadataRepository) delete(queryer Queryer, trackMetadataId int) (err error) {
	query := `
		DELETE FROM track_metadata
		WHERE track_metadata_id = :track_metadata_id
	`
	args := map[string]interface{}{
		"track_metadata_id": trackMetadataId,
	}
	_, err = queryer.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("id", trackMetadataId).Msg("Failed to delete track metadata")
		return err
	}

	log.Debug().Int("id", trackMetadataId).Msg("Track metadata deleted successfully")
	return nil
}

func (r TrackMetadataRepository) CountByAlbum(albumId int) (count int, err error) {
	log.Debug().Int("albumId", albumId).Msg("Counting track metadata by albumId")
	return r.countByAlbum(r.Db, albumId)
}

func (r TrackMetadataRepository) CountByAlbumTx(tx *sqlx.Tx, albumId int) (count int, err error) {
	log.Debug().Int("albumId", albumId).Msg("Counting track metadata by albumId transactional")
	return r.countByAlbum(tx, albumId)
}

func (r TrackMetadataRepository) countByAlbum(queryer Queryer, albumId int) (count int, err error) {
	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE album_id = :album_id
	`
	args := map[string]interface{}{
		"album_id": albumId,
	}
	err = queryer.Get(&count, query, args)
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to count tracks by albumId")
		return 0, err
	}

	log.Debug().Int("count", count).Int("albumId", albumId).Msg("Counted tracks by albumId successfully")
	return count, nil
}

func (r TrackMetadataRepository) CountByArtist(artistId int) (count int, err error) {
	log.Debug().Int("artistId", artistId).Msg("Counting track metadata by artistId")
	return r.countByArtist(r.Db, artistId)
}

func (r TrackMetadataRepository) CountByArtistTx(tx *sqlx.Tx, artistId int) (count int, err error) {
	log.Debug().Int("artistId", artistId).Msg("Counting track metadata by artistId transactional")
	return r.countByArtist(tx, artistId)
}

func (r TrackMetadataRepository) countByArtist(queryer Queryer, artistId int) (count int, err error) {
	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE artist_id = :artist_id
	`
	args := map[string]interface{}{
		"artist_id": artistId,
	}
	err = queryer.Get(&count, query, args)
	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to count tracks by artistId")
		return 0, err
	}

	log.Debug().Int("count", count).Int("artistId", artistId).Msg("Counted tracks by artistId successfully")
	return count, nil
}

func (r TrackMetadataRepository) CountByGenre(genreId int) (count int, err error) {
	log.Debug().Int("genreId", genreId).Msg("Counting track metadata by genreId")
	return r.countByGenre(r.Db, genreId)
}

func (r TrackMetadataRepository) CountByGenreTx(tx *sqlx.Tx, genreId int) (count int, err error) {
	log.Debug().Int("genreId", genreId).Msg("Counting track metadata by genreId transactional")
	return r.countByGenre(tx, genreId)
}

func (r TrackMetadataRepository) countByGenre(queryer Queryer, genreId int) (count int, err error) {
	const query = `
		SELECT COUNT(*)
		FROM track_metadata
		WHERE genre_id = :genre_id
	`
	args := map[string]interface{}{
		"genre_id": genreId,
	}
	err = queryer.Get(&count, query, args)
	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to count tracks by genreId")
		return 0, err
	}

	log.Debug().Int("count", count).Int("genreId", genreId).Msg("Counted tracks by genreId successfully")
	return count, nil
}

func (r TrackMetadataRepository) IsExistsByTrackId(trackId int) (exists bool, err error) {
	log.Debug().Int("trackId", trackId).Msg("Checking if track metadata exists by trackId")
	return r.isExistsByTrackId(r.Db, trackId)
}

func (r TrackMetadataRepository) IsExistsByTrackIdTx(tx *sqlx.Tx, trackId int) (exists bool, err error) {
	log.Debug().Int("trackId", trackId).Msg("Checking if track metadata exists by trackId transactional")
	return r.isExistsByTrackId(tx, trackId)
}

func (r TrackMetadataRepository) isExistsByTrackId(queryer Queryer, trackId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM track_metadata
			WHERE track_id = :track_id
		)
	`
	args := map[string]interface{}{
		"track_id": trackId,
	}
	row, err := queryer.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("trackId", trackId).Msg("Failed to execute query to check track metadata existence")
		return false, err
	}
	defer row.Close()

	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("trackId", trackId).Msg("Failed to scan result of artist_handler existence check")
			return false, err
		}
	}

	if exists {
		log.Debug().Int("trackId", trackId).Msg("Artist exists")
	} else {
		log.Debug().Int("trackId", trackId).Msg("No artist_handler found")
	}
	return exists, nil
}
