package album_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/handlers/types"
	"music-metadata/internal/models"
	"music-metadata/internal/service/track_metadata_service"
	"net/http"
	"strconv"
)

// readResponseTrack godoc
// @Description Response structure containing details about a track.
// @Property TrackMetadataId (integer) Unique identifier for the track metadata.
// @Property TrackId (integer) Unique identifier for the track.
// @Property Title (string, optional) Title of the track.
// @Property AlbumId (integer, optional) Identifier for the associated album.
// @Property ArtistId (integer, optional) Identifier for the associated artist.
// @Property GenreId (integer, optional) Identifier for the genre.
// @Property Year (integer, optional) Year of release.
// @Property TrackNumber (integer, optional) The track's position in the album.
// @Property DiscNumber (integer, optional) The disc number for multi-disc albums.
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

// readResponse godoc
// @Description Response structure containing detailed information about an album and its tracks.
// @Property AlbumId (integer) Unique identifier for the album.
// @Property Title (string) Title of the album.
// @Property CoverId (integer, optional) Identifier for the album's cover.
// @Property TracksCount (integer) Number of tracks in the album.
// @Property TrackMetadataList (array) List of track metadata for the album.
type readResponse struct {
	AlbumId           int                 `json:"albumId"`
	Title             string              `json:"title"`
	CoverId           *int                `json:"coverId,omitempty"`
	TracksCount       int                 `json:"tracksCount"`
	TrackMetadataList []readResponseTrack `json:"trackMetadataList"`
}

// Read godoc
// @Summary Get detailed information about an album and its tracks by album id
// @Tags Albums
// @Accept json
// @Produce json
// @Param albumId path integer true "Album Identifier"
// @Success 200 {object} readResponse
// @Failure 400 {object} types.Error "Invalid albumId format"
// @Failure 500 {object} types.Error "Failed to fetch album with details"
// @Router /albums/{albumId} [get]
func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching album with details")

	albumIdStr := c.Param("albumId")
	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		log.Error().Err(err).Str("albumIdStr", albumIdStr).Msg("Invalid albumId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid albumId format",
		})
		return
	}
	log.Debug().Int("albumId", albumId).Msg("Url parameter read successfully")

	var album models.Album
	var trackMetadataList []track_metadata_service.TrackMetadataReadByAlbumId
	var coverId *int

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		album, err = h.AlbumService.Read(tx, albumId)
		if err != nil {
			return err
		}

		trackMetadataList, err = h.TrackMetadataService.ReadByAlbumId(tx, album.AlbumId)
		if err != nil {
			return err
		}

		coverId, err = h.CoverService.GetMostCommonCoverIdByAlbumId(tx, albumId)
		if err != nil {
			coverId = nil
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
		CoverId:           coverId,
		TracksCount:       len(trackMetadataList),
		TrackMetadataList: trackMetadataListResponse,
	}

	log.Debug().Int("albumId", albumId).Int("tracksCount", len(trackMetadataList)).Msg("Album fetched successfully")
	c.JSON(http.StatusOK, response)
}
