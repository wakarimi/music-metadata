package song_service

import (
	"bytes"
	"github.com/dhowden/tag"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/model"
	"strings"
)

func (s *Service) SongByAudioFileWithoutSha(tx *sqlx.Tx, audioFileId int) (song model.Song, err error) {
	file, err := s.AudioFileClient.Download(audioFileId)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to download audio file")
		return model.Song{}, err
	}

	metadata, err := extractMetadata(file)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to extract file's metadata")
		return model.Song{}, err
	}

	albumId, err := s.getOrCreateAlbum(tx, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get album")
		return model.Song{}, err
	}
	artistId, err := s.getOrCreateArtist(tx, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artist")
		return model.Song{}, err
	}
	genreId, err := s.getOrCreateGenre(tx, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get genre")
		return model.Song{}, err
	}

	song = model.Song{
		AudioFileId: audioFileId,
		Title:       getTitle(metadata),
		AlbumId:     albumId,
		ArtistId:    artistId,
		GenreId:     genreId,
		Year:        getYear(metadata),
		SongNumber:  getSongNumber(metadata),
		DiscNumber:  getDiscNumber(metadata),
		Lyrics:      getLyrics(metadata),
	}

	return song, nil
}

func (s *Service) getOrCreateAlbum(tx *sqlx.Tx, metadata tag.Metadata) (albumId *int, err error) {
	title := strings.TrimSpace(metadata.Album())
	if len(title) == 0 {
		return nil, nil
	}

	exists, err := s.AlbumService.IsExistsByTitle(tx, title)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to check album existence")
		return nil, err
	}

	if exists {
		album, err := s.AlbumService.GetByTitle(tx, title)
		if err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to get album")
			return nil, err
		}
		return &album.AlbumId, nil
	} else {
		album, err := s.AlbumService.Create(tx, model.Album{
			Title: title,
		})
		if err != nil {
			log.Error().Err(err).Str("title", title).Msg("Failed to create album")
			return nil, err
		}
		return &album.AlbumId, nil
	}
}

func (s *Service) getOrCreateArtist(tx *sqlx.Tx, metadata tag.Metadata) (artistId *int, err error) {
	name := strings.TrimSpace(metadata.Artist())
	if len(name) == 0 {
		return nil, nil
	}

	exists, err := s.ArtistService.IsExistsByName(tx, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to check artist existence")
		return nil, err
	}

	if exists {
		artist, err := s.ArtistService.GetByName(tx, name)
		if err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to get artist")
			return nil, err
		}
		return &artist.ArtistId, nil
	} else {
		artist, err := s.ArtistService.Create(tx, model.Artist{
			Name: name,
		})
		if err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to create artist")
			return nil, err
		}
		return &artist.ArtistId, nil
	}
}

func (s *Service) getOrCreateGenre(tx *sqlx.Tx, metadata tag.Metadata) (genreId *int, err error) {
	name := strings.TrimSpace(metadata.Genre())
	if len(name) == 0 {
		return nil, nil
	}

	exists, err := s.GenreService.IsExistsByName(tx, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to check genre existence")
		return nil, err
	}

	if exists {
		genre, err := s.GenreService.GetByName(tx, name)
		if err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to get genre")
			return nil, err
		}
		return &genre.GenreId, nil
	} else {
		genre, err := s.GenreService.Create(tx, model.Genre{
			Name: name,
		})
		if err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to create genre")
			return nil, err
		}
		return &genre.GenreId, nil
	}
}

func extractMetadata(trackData []byte) (metadata tag.Metadata, err error) {
	r := bytes.NewReader(trackData)
	return tag.ReadFrom(r)
}

func getTitle(metadata tag.Metadata) *string {
	title := strings.TrimSpace(metadata.Title())
	if len(title) == 0 {
		return nil
	}
	return &title
}

func getYear(metadata tag.Metadata) *int {
	year := metadata.Year()
	if year == 0 {
		return nil
	}
	return &year
}

func getSongNumber(metadata tag.Metadata) *int {
	trackNumber, _ := metadata.Track()
	if trackNumber == 0 {
		return nil
	}
	return &trackNumber
}

func getDiscNumber(metadata tag.Metadata) *int {
	discNumber, _ := metadata.Disc()
	if discNumber == 0 {
		return nil
	}
	return &discNumber
}

func getLyrics(metadata tag.Metadata) *string {
	lyrics := strings.TrimSpace(metadata.Lyrics())
	if len(lyrics) == 0 {
		return nil
	}
	return &lyrics
}
