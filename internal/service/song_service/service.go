package song_service

import (
	"music-metadata/internal/client/music_files_client/audio_file_client"
	"music-metadata/internal/database/repository/song_repo"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/genre_service"
)

type Service struct {
	SongRepo song_repo.Repo

	AlbumService  album_service.Service
	ArtistService artist_service.Service
	GenreService  genre_service.Service

	AudioFileClient audio_file_client.Client
}

func NewService(songRepo song_repo.Repo,
	albumService album_service.Service,
	artistService artist_service.Service,
	genreService genre_service.Service,
	audioFileClient audio_file_client.Client) (s *Service) {

	s = &Service{
		SongRepo:        songRepo,
		AlbumService:    albumService,
		ArtistService:   artistService,
		GenreService:    genreService,
		AudioFileClient: audioFileClient,
	}

	return s
}
