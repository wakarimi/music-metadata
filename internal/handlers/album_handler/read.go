package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/track_metadata_service"
	"net/http"
	"strconv"
)

type readResponseTrack struct {
	TrackMetadataId int     `json:"trackMetadataId"`
	TrackId         int     `json:"trackId"`
	Title           *string `json:"title"`
	AlbumId         *int    `json:"albumId"`
	ArtistId        *int    `json:"artistId"`
	GenreId         *int    `json:"genreId"`
	Year            *int    `json:"year"`
	TrackNumber     *int    `json:"trackNumber"`
	DiscNumber      *int    `json:"discNumber"`
}

type readResponse struct {
	AlbumId           int                 `json:"albumId"`
	Title             string              `json:"title"`
	CoverId           *int                `json:"coverId,omitempty"`
	TracksCount       int                 `json:"tracksCount"`
	TrackMetadataList []readResponseTrack `json:"trackMetadataList"`
}

// Read godoc
// @Summary Get detailed information about an album and its tracks by album id
// @Tags albums
// @Accept json
// @Produce json
// @Param albumId path integer true "Album Identifier"
// @Success 200 {object} readResponse
// @Failure 400 {object} types.Error "Invalid albumId format"
// @Failure 500 {object} types.Error "Failed to fetch album with details"
// @Router /albums/{albumId} [get]
func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching album with ditails")

	albumIdStr := c.Param("albumId")
	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid albumId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid albumId format",
		})
		return
	}
	log.Debug().Int("albumId", albumId).Msg("Url parameter read successfully")

	var album album_service.AlbumRead
	var trackMetadataList []track_metadata_service.TrackMetadataReadByAlbumId

	h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		album, err = h.AlbumService.Read(tx, albumId)
		if err != nil {
			log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album")
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to fetch album",
			})
			return err
		}

		trackMetadataList, err = h.TrackMetadataService.ReadByAlbumId(tx, album.AlbumId)
		if err != nil {
			log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album's tracks")
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to fetch album's tracks",
			})
			return err
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Int("albumId", albumId).Msg("Failed to fetch album with details")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch album with details",
		})
		return
	}

	trackMetadataListResponse := make([]readResponseTrack, len(trackMetadataList))
	for i, trackMetadata := range trackMetadataList {
		trackMetadataListResponse[i] = readResponseTrack{
			TrackMetadataId: trackMetadata.TrackMetadataId,
			TrackId:         trackMetadata.TrackMetadataId,
			Title:           trackMetadata.Title,
			AlbumId:         trackMetadata.AlbumId,
			ArtistId:        trackMetadata.ArtistId,
			GenreId:         trackMetadata.GenreId,
			Year:            trackMetadata.Year,
			TrackNumber:     trackMetadata.TrackNumber,
			DiscNumber:      trackMetadata.DiscNumber,
		}
	}

	response := readResponse{
		AlbumId:           album.AlbumId,
		Title:             album.Title,
		CoverId:           album.CoverId,
		TracksCount:       album.TracksCount,
		TrackMetadataList: trackMetadataListResponse,
	}

	log.Debug().Int("albumId", albumId).Msg("Album fetched successfully")
	c.JSON(http.StatusOK, response)
}
