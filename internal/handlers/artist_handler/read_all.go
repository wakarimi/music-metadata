package artist_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/models"
	"net/http"
)

// readAllResponseItem godoc
// @Description Artist details
// @Property ArtistId (integer) The unique identifier of the artist
// @Property Name (string) Name of the artist
// @Property CoverId (integer, optional) Identifier of the artist's most popular cover
// @Property TracksCount (integer) Number of tracks by the artist
type readAllResponseItem struct {
	ArtistId          int    `json:"artistId"`
	Name              string `json:"name"`
	MostPopularCovers []int  `json:"mostPopularCoversIds,omitempty"`
	TracksCount       int    `json:"tracksCount"`
}

// readAllResponse godoc
// @Description List of all artists
// @Property Artists (array) List of artist details
type readAllResponse struct {
	Artists []readAllResponseItem `json:"artists"`
}

// ReadAll godoc
// @Summary Get all artists
// @Tags Artists
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse
// @Failure 500 {object} types.Error "Failed to fetch all artists"
// @Router /artists [get]
func (h *Handler) ReadAll(c *gin.Context) {
	log.Debug().Msg("Fetching all artists")

	var artists []models.Artist
	var coverIds [][]int
	var tracksCounts []int

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artists, err = h.ArtistService.ReadAll(tx)
		if err != nil {
			return err
		}

		for _, artist := range artists {
			coverId, err := h.CoverService.GetMostCommonCoverIdsByArtistId(tx, artist.ArtistId, 4)
			if err != nil {
				coverId = make([]int, 0)
			}
			coverIds = append(coverIds, coverId)

			var tracksCount int
			trackMetadataList, err := h.TrackMetadataService.ReadByArtistId(tx, artist.ArtistId)
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
		log.Error().Err(err).Msg("Failed to fetch all artists")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch all artists",
		})
		return
	}

	artistsResponse := make([]readAllResponseItem, len(artists))
	for i, artist := range artists {
		artistResponse := readAllResponseItem{
			ArtistId:          artist.ArtistId,
			Name:              artist.Name,
			MostPopularCovers: coverIds[i],
			TracksCount:       tracksCounts[i],
		}
		artistsResponse[i] = artistResponse
	}

	log.Debug().Int("count", len(artistsResponse)).Msg("All artists fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Artists: artistsResponse,
	})
}
