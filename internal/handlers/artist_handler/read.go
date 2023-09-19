package artist_handler

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
// @Property ArtistId (integer, optional) Identifier for the associated artist_handler.
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

// readResponseArtist godoc
// @Description Response structure containing detailed information about an artist and his tracks.
// @Property ArtistId (integer) Unique identifier for the artist.
// @Property Name (string) Name of the artist.
// @Property MostPopularCoverId (integer, optional) Identifier for the artist's most popular cover.
// @Property TracksCount (integer) Number of tracks performed by the artist.
// @Property TrackMetadataList (array) List of track metadata for the artist.
type readResponseArtist struct {
	ArtistId            int                 `json:"artistId"`
	Name                string              `json:"name"`
	MostPopularCoverIds []int               `json:"mostPopularCoverIds,omitempty"`
	TracksCount         int                 `json:"tracksCount"`
	TrackMetadataList   []readResponseTrack `json:"trackMetadataList"`
}

// ReadArtist godoc
// @Summary Get detailed information about an artist and his tracks by artist id
// @Tags Artists
// @Accept json
// @Produce json
// @Param artistId path integer true "Artist Identifier"
// @Success 200 {object} readResponseArtist
// @Failure 400 {object} types.Error "Invalid artistId format"
// @Failure 500 {object} types.Error "Failed to fetch artist with details"
// @Router /artists/{artistId} [get]
func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching artist with details")

	artistIdStr := c.Param("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil {
		log.Error().Err(err).Str("artistIdStr", artistIdStr).Msg("Invalid artistId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid artistId format",
		})
		return
	}

	var artist models.Artist
	var trackMetadataList []track_metadata_service.TrackMetadataReadByArtistId
	var mostPopularCoverIds []int

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		artist, err = h.ArtistService.Read(tx, artistId)
		if err != nil {
			return err
		}

		trackMetadataList, err = h.TrackMetadataService.ReadByArtistId(tx, artist.ArtistId)
		if err != nil {
			return err
		}

		coverIds, err := h.CoverService.GetMostCommonCoverIdsByArtistId(tx, artistId, 4)
		if err != nil {
			mostPopularCoverIds = make([]int, 0)
		} else {
			mostPopularCoverIds = coverIds
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Int("artistId", artistId).Msg("Failed to fetch artist with details")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch artist with details",
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

	response := readResponseArtist{
		ArtistId:            artist.ArtistId,
		Name:                artist.Name,
		MostPopularCoverIds: mostPopularCoverIds,
		TracksCount:         len(trackMetadataList),
		TrackMetadataList:   trackMetadataListResponse,
	}

	log.Debug().Int("artistId", artistId).Int("tracksCount", len(trackMetadataList)).Msg("Artist fetched successfully")
	c.JSON(http.StatusOK, response)
}
