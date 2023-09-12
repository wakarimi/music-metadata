package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-metadata/internal/clients/musicfilesclient"
	"music-metadata/internal/clients/musicfilesclient/trackrequests"
	"music-metadata/internal/context"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers/album_handler"
	"music-metadata/internal/handlers/artist"
	"music-metadata/internal/handlers/genre"
	"music-metadata/internal/handlers/track_metadata"
	"music-metadata/internal/middleware"
	"music-metadata/internal/service"
	"music-metadata/internal/service/albumservice"
)

func SetupRouter(ctx *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	musicFilesClient := musicfilesclient.NewClient(ctx.Config.HttpServer.MusicFilesAddress)
	trackClient := trackrequests.NewTrackClient(musicFilesClient)

	txManager := service.NewTransactionManager(*ctx.Db)

	albumRepo := repository.NewAlbumRepository(ctx.Db)
	artistRepo := repository.NewArtistRepository(ctx.Db)
	genreRepo := repository.NewGenreRepository(ctx.Db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(ctx.Db)

	albumService := albumservice.NewService(txManager, albumRepo, trackMetadataRepo, trackClient)

	albumHandler := album_handler.NewHandler(*albumService)
	artistHandler := artist.NewArtistHandler(artistRepo)
	genreHandler := genre.NewGenreHandler(genreRepo)
	musicHandler := track_metadata.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-metadata-service")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		albums := api.Group("/albums")
		{
			albums.GET("/", func(c *gin.Context) { albumHandler.ReadAll(c, ctx) })
		}
		artists := api.Group("/artists")
		{
			artists.GET("/", artistHandler.GetAll)
		}
		genres := api.Group("/genres")
		{
			genres.GET("/", genreHandler.GetAll)
		}

		api.POST("/scan", func(c *gin.Context) { musicHandler.Scan(c, &ctx.Config.HttpServer) })
	}

	log.Debug().Msg("Router setup successfully")
	return r
}
