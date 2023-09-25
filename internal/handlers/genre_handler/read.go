package genre_handler

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

// readResponseGenre godoc
// @Description Response structure containing detailed information about an genre and his tracks.
// @Property GenreId (integer) Unique identifier for the genre.
// @Property Name (string) Name of the genre.
// @Property MostPopularCoverIds (array, optional) Identifier for the genre's most popular cover.
// @Property TracksCount (integer) Number of tracks performed by the genre.
// @Property TrackMetadataList (array) List of track metadata for the genre.
type readResponseGenre struct {
	GenreId             int                 `json:"genreId"`
	Name                string              `json:"name"`
	MostPopularCoverIds []int               `json:"mostPopularCoverIds,omitempty"`
	TracksCount         int                 `json:"tracksCount"`
	TrackMetadataList   []readResponseTrack `json:"trackMetadataList"`
}

// ReadGenre godoc
// @Summary Get detailed information about an genre and his tracks by genre id
// @Tags Genres
// @Accept json
// @Produce json
// @Param genreId path integer true "Genre Identifier"
// @Success 200 {object} readResponseGenre
// @Failure 400 {object} types.Error "Invalid genreId format"
// @Failure 500 {object} types.Error "Failed to fetch genre with details"
// @Router /genres/{genreId} [get]
func (h *Handler) Read(c *gin.Context) {
	log.Debug().Msg("Fetching genre with details")

	genreIdStr := c.Param("genreId")
	genreId, err := strconv.Atoi(genreIdStr)
	if err != nil {
		log.Error().Err(err).Str("genreIdStr", genreIdStr).Msg("Invalid genreId format")
		c.JSON(http.StatusBadRequest, types.Error{
			Error: "Invalid genreId format",
		})
		return
	}

	var genre models.Genre
	var trackMetadataList []track_metadata_service.TrackMetadataReadByGenreId
	var mostPopularCoverIds []int

	err = h.TransactionManager.WithTransaction(func(tx *sqlx.Tx) (err error) {
		genre, err = h.GenreService.Read(tx, genreId)
		if err != nil {
			return err
		}

		trackMetadataList, err = h.TrackMetadataService.ReadByGenreId(tx, genre.GenreId)
		if err != nil {
			return err
		}

		coverIds, err := h.CoverService.GetMostCommonCoverIdsByGenreId(tx, genreId, 4)
		if err != nil {
			mostPopularCoverIds = make([]int, 0)
		} else {
			mostPopularCoverIds = coverIds
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Int("genreId", genreId).Msg("Failed to fetch genre with details")
		c.JSON(http.StatusInternalServerError, types.Error{
			Error: "Failed to fetch genre with details",
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

	response := readResponseGenre{
		GenreId:             genre.GenreId,
		Name:                genre.Name,
		MostPopularCoverIds: mostPopularCoverIds,
		TracksCount:         len(trackMetadataList),
		TrackMetadataList:   trackMetadataListResponse,
	}

	log.Debug().Int("genreId", genreId).Int("tracksCount", len(trackMetadataList)).Msg("Genre fetched successfully")
	c.JSON(http.StatusOK, response)
}
