package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/config"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers/album"
	"music-metadata/internal/handlers/artist"
	"music-metadata/internal/handlers/genre"
	"music-metadata/internal/handlers/track_metadata"
	"music-metadata/internal/middleware"
)

func SetupRouter(httpServerConfig *config.HttpServer, db *sqlx.DB) *gin.Engine {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	albumRepo := repository.NewAlbumRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(db)

	albumHandler := album.NewAlbumHandler(albumRepo)
	artistHandler := artist.NewArtistHandler(artistRepo)
	genreHandler := genre.NewGenreHandler(genreRepo)
	musicHandler := track_metadata.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r := gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-metadata-service")
	{
		albums := api.Group("/albums")
		{
			albums.GET("/", albumHandler.GetAll)
		}
		artists := api.Group("/artists")
		{
			artists.GET("/", artistHandler.GetAll)
		}
		genres := api.Group("/genres")
		{
			genres.GET("/", genreHandler.GetAll)
		}

		api.POST("/scan", func(c *gin.Context) { musicHandler.Scan(c, httpServerConfig) })
	}

	return r
}
