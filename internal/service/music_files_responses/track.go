package music_files_responses

type TrackGet struct {
	TrackId   int    `json:"trackId"`
	CoverId   *int   `json:"coverId"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
}
