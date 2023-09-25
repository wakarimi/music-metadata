package genre_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/models"
	"net/http"
)

// readAllResponseItem godoc
// @Description Genre details
// @Property GenreId (integer) The unique identifier of the genre
// @Property Name (string) Name of the genre
// @Property CoverId (integer, optional) Identifier of the genre's most popular cover
// @Property TracksCount (integer) Number of tracks by the genre
type readAllResponseItem struct {
	GenreId           int    `json:"genreId"`
	Name              string `json:"name"`
	MostPopularCovers []int  `json:"mostPopularCoversIds,omitempty"`
	TracksCount       int    `json:"tracksCount"`
}

// readAllResponse godoc
// @Description List of all genres
// @Property Genres (array) List of genre details
type readAllResponse struct {
	Genres []readAllResponseItem `json:"genres"`
}

// ReadAll godoc
// @Summary Get all genres
// @Tags Genres
// @Accept json
// @Produce json
// @Success 200 {object} readAllResponse
// @Failure 500 {object} types.Error "Failed to fetch all genres"
// @Router /genres [get]
func (h *Handler) ReadAll(c *gin.Context) {
	log.Debug().Msg("Fetching all genres")

	var genres []models.Genre
	var coverIds [][]int
	var tracksCounts []int

	err := h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genres, err = h.GenreService.ReadAll(tx)
		if err != nil {
			return err
		}

		for _, genre := range genres {
			coverId, err := h.CoverService.GetMostCommonCoverIdsByGenreId(tx, genre.GenreId, 4)
			if err != nil {
				coverId = make([]int, 0)
			}
			coverIds = append(coverIds, coverId)

			var tracksCount int
			trackMetadataList, err := h.TrackMetadataService.ReadByGenreId(tx, genre.GenreId)
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
		log.Error().Err(err).Msg("Failed to fetch all genres")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch all genres",
		})
		return
	}

	genresResponse := make([]readAllResponseItem, len(genres))
	for i, genre := range genres {
		genreResponse := readAllResponseItem{
			GenreId:           genre.GenreId,
			Name:              genre.Name,
			MostPopularCovers: coverIds[i],
			TracksCount:       tracksCounts[i],
		}
		genresResponse[i] = genreResponse
	}

	log.Debug().Int("count", len(genresResponse)).Msg("All genres fetched successfully")
	c.JSON(http.StatusOK, readAllResponse{
		Genres: genresResponse,
	})
}
