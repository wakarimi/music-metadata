package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/config"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers"
	"music-metadata/internal/middleware"
)

func SetupRouter(httpServerConfig *config.HttpServer, db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	albumRepo := repository.NewAlbumRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(db)

	albumHandler := handlers.NewAlbumHandler(albumRepo)
	musicHandler := handlers.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r := gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-metadata-service")
	{
		albums := api.Group("/albums")
		{
			albums.GET("/", func(c *gin.Context) { albumHandler.GetAll(c) })
		}

		api.POST("/scan", func(c *gin.Context) { musicHandler.Scan(c, httpServerConfig) })
	}

	return r
}
