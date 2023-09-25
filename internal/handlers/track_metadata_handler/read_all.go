package track_metadata_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"net/http"
)

// readAllResponseAlbum godoc
type readAllResponseAlbum struct {
	AlbumId int    `json:"albumId"`
	Title   string `json:"title"`
}

// readAllResponseArtist godoc
type readAllResponseArtist struct {
	ArtistId int    `json:"artistId"`
	Name     string `json:"name"`
}

// readAllResponseGenre godoc
type readAllResponseGenre struct {
	GenreId int    `json:"genreId"`
	Name    string `json:"name"`
}

// readAllResponseItem godoc
type readAllResponseItem struct {
	TrackMetadataId int                    `json:"trackMetadataId"`
	TrackId         int                    `json:"trackId"`
	Title           *string                `json:"title"`
	Album           *readAllResponseAlbum  `json:"album,omitempty"`
	Artist          *readAllResponseArtist `json:"artist,omitempty"`
	Genre           *readAllResponseGenre  `json:"genre,omitempty"`
	Year            *int                   `json:"year"`
	TrackNumber     *int                   `json:"trackNumber"`
	DiscNumber      *int                   `json:"discNumber"`
}

// readAllResponse godoc
type readAllResponse struct {
	Genres []readAllResponseItem `json:"trackMetadataList"`
}

// ReadAll godoc
// @Summary Get all track metadata
// @Tags TrackMetadata
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse
// @Failure 500 {object} types.Error "Failed to fetch all genres"
// @Router /track-metadata [get]
func (h *Handler) ReadAll(c *gin.Context) {
	log.Debug().Msg("Fetching all track metadata")

	var trackMetadataListResponse []readAllResponseItem

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		trackMetadataList, err := h.TrackMetadataService.ReadAll(tx)
		if err != nil {
			return err
		}

		for _, trackMetadata := range trackMetadataList {
			var album *readAllResponseAlbum
			var artist *readAllResponseArtist
			var genre *readAllResponseGenre

			if trackMetadata.AlbumId != nil {
				albumModel, err := h.AlbumService.Read(tx, *trackMetadata.AlbumId)
				if err == nil {
					album = &readAllResponseAlbum{
						AlbumId: albumModel.AlbumId,
						Title:   albumModel.Title,
					}
				}
			}

			if trackMetadata.ArtistId != nil {
				artistModel, err := h.ArtistService.Read(tx, *trackMetadata.ArtistId)
				if err == nil {
					artist = &readAllResponseArtist{
						ArtistId: artistModel.ArtistId,
						Name:     artistModel.Name,
					}
				}
			}

			if trackMetadata.GenreId != nil {
				genreModel, err := h.GenreService.Read(tx, *trackMetadata.GenreId)
				if err == nil {
					genre = &readAllResponseGenre{
						GenreId: genreModel.GenreId,
						Name:    genreModel.Name,
					}
				}
			}

			newTrackMetadataResponse := readAllResponseItem{
				TrackMetadataId: trackMetadata.TrackMetadataId,
				TrackId:         trackMetadata.TrackMetadataId,
				Title:           trackMetadata.Title,
				Album:           album,
				Artist:          artist,
				Genre:           genre,
				Year:            trackMetadata.Year,
				TrackNumber:     trackMetadata.TrackNumber,
				DiscNumber:      trackMetadata.DiscNumber,
			}

			trackMetadataListResponse = append(trackMetadataListResponse, newTrackMetadataResponse)
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all track metadata")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch all track metadata",
		})
		return
	}

	log.Debug().Int("count", len(trackMetadataListResponse)).Msg("All track metadata fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Genres: trackMetadataListResponse,
	})
}
