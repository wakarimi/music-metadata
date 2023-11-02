package artist_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"modernc.org/mathutil"
	"music-metadata/internal/handlers/response"
	"music-metadata/internal/model"
	"net/http"
	"strconv"
)

// getAllResponseItem represents a single artist item in the GetAllArtists API response.
type getAllResponseItem struct {
	// Unique identifier for the artist.
	ArtistId int `json:"artistId"`
	// Name of the artist.
	Name string `json:"name"`
	// Optional array of best cover IDs for the artist.
	BestCover *[]int `json:"bestCovers,omitempty"`
}

// getAllResponse represents the response model for GetAllArtists API.
type getAllResponse struct {
	// Array of artists.
	Artists []getAllResponseItem `json:"artists"`
}

// GetAll retrieves a list of all artists with optional best covers.
// @Summary Retrieve all artists
// @Description Retrieves a list of all artists, including their best covers if requested.
// @Tags Artists
// @Accept  json
// @Produce  json
// @Param   bestCovers   query   int     false       "Number of best covers for each artist to retrieve"
// @Success 200 {object} getAllResponse "Success response with a list of artists and optional best covers for each"
// @Failure 400 {object} response.Error "Invalid bestCovers format"
// @Failure 500 {object} response.Error "Internal Server Error"
// @Router /artists [get]
func (h *Handler) GetAll(c *gin.Context) {
	log.Debug().Msg("Getting artists")

	bestCoversStr := c.DefaultQuery("bestCovers", "0")
	bestCoversInt, err := strconv.Atoi(bestCoversStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bestCovers format")
		c.JSON(http.StatusBadRequest, response.Error{
			Message: "Invalid bestCovers format",
			Reason:  err.Error(),
		})
		return
	}

	var artists []model.Artist
	var bestCovers [][]int
	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artists, err = h.ArtistService.GetAll(tx)
		if err != nil {
			return err
		}
		if bestCoversInt > 0 {
			for _, artist := range artists {
				bestCoversItem, err := h.CoverService.CalcBestCoversForArtist(tx, artist.ArtistId)
				if err != nil {
					return err
				}
				bestCovers = append(bestCovers, bestCoversItem)
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get artists")
		c.JSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to get artists",
			Reason:  err.Error(),
		})
		return
	}

	artistsResponseItems := make([]getAllResponseItem, len(artists))
	for i, artist := range artists {
		bestCoversForArtist := make([]int, 0)
		if bestCoversInt > 0 {
			for j := 0; j < mathutil.Min(bestCoversInt, len(bestCovers[i])); j++ {
				bestCoversForArtist = append(bestCoversForArtist, bestCovers[i][j])
			}
		}
		var bestCoverForArtistResponse *[]int
		if len(bestCoversForArtist) > 0 {
			bestCoverForArtistResponse = &bestCoversForArtist
		} else {
			bestCoverForArtistResponse = nil
		}
		artistsResponseItems[i] = getAllResponseItem{
			ArtistId:  artist.ArtistId,
			Name:      artist.Name,
			BestCover: bestCoverForArtistResponse,
		}
	}

	log.Debug().Msg("Artists got successfully")
	c.JSON(http.StatusOK, getAllResponse{
		Artists: artistsResponseItems,
	})
}
