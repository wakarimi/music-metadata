package cover_service

import (
	"music-metadata/internal/client/music_files_client/audio_file_client"
	"music-metadata/internal/service/song_service"
)

type Service struct {
	SongService song_service.Service

	AudioFileClient audio_file_client.Client
}

func NewService(songService song_service.Service,
	audioFileClient audio_file_client.Client) (s *Service) {

	s = &Service{
		SongService:     songService,
		AudioFileClient: audioFileClient,
	}

	return s
}
