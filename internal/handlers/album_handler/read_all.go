package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/models"
	"net/http"
)

// readAllResponseItem godoc
// @Description Album details
// @Property AlbumId (integer) The unique identifier of the album
// @Property Title (string) Title of the album
// @Property CoverId (integer, optional) Identifier of the album cover
// @Property TracksCount (integer) Number of tracks in the album
type readAllResponseItem struct {
	AlbumId     int    `json:"albumId"`
	Title       string `json:"title"`
	CoverId     *int   `json:"coverId,omitempty"`
	TracksCount int    `json:"tracksCount"`
}

// readAllResponse godoc
// @Description List of all albums
// @Property Albums (array) List of album details
type readAllResponse struct {
	Albums []readAllResponseItem `json:"albums"`
}

// ReadAll godoc
// @Summary Get all albums
// @Tags Albums
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse
// @Failure 500 {object} types.Error "Failed to fetch all album"
// @Router /albums [get]
func (h *Handler) ReadAll(c *gin.Context) {
	log.Debug().Msg("Fetching all albums")

	var albums []models.Album
	var coverIds []*int
	var tracksCounts []int

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		albums, err = h.AlbumService.ReadAll(tx)
		if err != nil {
			return err
		}

		for _, album := range albums {
			coverId, err := h.CoverService.GetMostCommonCoverIdByAlbumId(tx, album.AlbumId)
			if err != nil {
				coverId = nil
			}
			coverIds = append(coverIds, coverId)

			var tracksCount int
			trackMetadataList, err := h.TrackMetadataService.ReadByAlbumId(tx, album.AlbumId)
			if err != nil {
				tracksCount = 0
			} else {
				tracksCount = len(trackMetadataList)
			}

			tracksCounts = append(tracksCounts, tracksCount)
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all albums")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch all album",
		})
		return
	}

	albumsResponse := make([]readAllResponseItem, len(albums))
	for i, album := range albums {
		albumResponse := readAllResponseItem{
			AlbumId:     album.AlbumId,
			Title:       album.Title,
			CoverId:     coverIds[i],
			TracksCount: tracksCounts[i],
		}
		albumsResponse[i] = albumResponse
	}

	log.Debug().Int("count", len(albumsResponse)).Msg("All albums fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Albums: albumsResponse,
	})
}
