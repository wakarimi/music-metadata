package trackrequests

import (
	"music-metadata/internal/clients/musicfilesclient"
)

type TrackClient struct {
	musicFilesClient *musicfilesclient.Client
}

func NewTrackClient(musicFilesClient *musicfilesclient.Client) (trackClient TrackClient) {
	trackClient = TrackClient{
		musicFilesClient: musicFilesClient,
	}

	return trackClient
}
