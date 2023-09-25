package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-metadata/internal/clients/music_files_client"
	"music-metadata/internal/clients/music_files_client/track_requests"
	"music-metadata/internal/context"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers/album_handler"
	"music-metadata/internal/handlers/artist_handler"
	"music-metadata/internal/handlers/genre_handler"
	"music-metadata/internal/handlers/track_metadata"
	"music-metadata/internal/handlers/track_metadata_handler"
	"music-metadata/internal/middleware"
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/cover_service"
	"music-metadata/internal/service/genre_service"
	"music-metadata/internal/service/track_metadata_service"
)

func SetupRouter(appCtx *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	musicFilesClient := music_files_client.NewClient(appCtx.Config.HttpServer.MusicFilesAddress)
	trackClient := track_requests.NewTrackClient(musicFilesClient)

	txManager := service.NewTransactionManager(*appCtx.Db)

	albumRepo := repository.NewAlbumRepository(appCtx.Db)
	artistRepo := repository.NewArtistRepository(appCtx.Db)
	genreRepo := repository.NewGenreRepository(appCtx.Db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(appCtx.Db)

	albumService := album_service.NewService(albumRepo, trackMetadataRepo, trackClient)
	artistService := artist_service.NewService(artistRepo, trackMetadataRepo, trackClient)
	genreService := genre_service.NewService(genreRepo, trackMetadataRepo, trackClient)
	trackMetadataService := track_metadata_service.NewService(txManager, trackMetadataRepo)
	coverService := cover_service.NewService(albumRepo, artistRepo, genreRepo, trackMetadataRepo, trackClient)

	trackMetadataHandler := track_metadata_handler.NewHandler(txManager, *trackMetadataService, *albumService, *artistService, *genreService)
	albumHandler := album_handler.NewHandler(txManager, *albumService, *trackMetadataService, *coverService)
	artistHandler := artist_handler.NewHandler(txManager, *artistService, *trackMetadataService, *coverService)
	genreHandler := genre_handler.NewHandler(txManager, *genreService, *trackMetadataService, *coverService)
	musicHandler := track_metadata.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	api := r.Group("/api/music-metadata-service")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		trackMetadata := api.Group("/track-metadata")
		{
			trackMetadata.GET("/", trackMetadataHandler.ReadAll)
		}
		albums := api.Group("/albums")
		{
			albums.GET("/", albumHandler.ReadAll)
			albums.GET("/:albumId", albumHandler.Read)
		}
		artists := api.Group("/artists")
		{
			artists.GET("/", artistHandler.ReadAll)
			artists.GET("/:artistId", artistHandler.Read)
		}
		genres := api.Group("/genres")
		{
			genres.GET("/", genreHandler.ReadAll)
			genres.GET("/:genreId", genreHandler.Read)
		}

		api.POST("/scan", func(c *gin.Context) { musicHandler.Scan(c, &appCtx.Config.HttpServer) })
	}

	log.Debug().Msg("Router setup successfully")
	return r
}
