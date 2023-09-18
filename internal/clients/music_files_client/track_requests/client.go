package track_requests

import (
	"music-metadata/internal/clients/music_files_client"
)

type TrackClient struct {
	musicFilesClient *music_files_client.Client
}

func NewTrackClient(musicFilesClient *music_files_client.Client) (trackClient TrackClient) {
	trackClient = TrackClient{
		musicFilesClient: musicFilesClient,
	}

	return trackClient
}
