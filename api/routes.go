package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/config"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers/album_handler"
	"music-metadata/internal/handlers/artist"
	"music-metadata/internal/handlers/genre"
	"music-metadata/internal/handlers/track_metadata"
	"music-metadata/internal/middleware"
	"music-metadata/internal/service/album_service"
)

func SetupRouter(httpServerConfig *config.HttpServer, db *sqlx.DB) *gin.Engine {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	albumRepo := repository.NewAlbumRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(db)

	albumService := album_service.NewService(albumRepo)

	albumHandler := album_handler.NewHandler(*albumService)
	artistHandler := artist.NewArtistHandler(artistRepo)
	genreHandler := genre.NewGenreHandler(genreRepo)
	musicHandler := track_metadata.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r := gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-metadata-service")
	{
		albums := api.Group("/albums")
		{
			albums.GET("/", albumHandler.ReadAll)
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
