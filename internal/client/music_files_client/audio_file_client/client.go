package audio_file_client

import (
	"music-metadata/internal/client/music_files_client"
)

type Client struct {
	audioFileClient *music_files_client.Client
}

func NewAudioFileClient(audioFileClient *music_files_client.Client) (trackClient Client) {
	trackClient = Client{
		audioFileClient: audioFileClient,
	}

	return trackClient
}
