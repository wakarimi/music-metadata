package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-metadata/internal/client/music_files_client"
	"music-metadata/internal/client/music_files_client/audio_file_client"
	"music-metadata/internal/context"
	"music-metadata/internal/database/repository/album_repo"
	"music-metadata/internal/database/repository/artist_repo"
	"music-metadata/internal/database/repository/genre_repo"
	"music-metadata/internal/database/repository/song_repo"
	"music-metadata/internal/handlers/album_handler"
	"music-metadata/internal/handlers/artist_handler"
	"music-metadata/internal/handlers/genre_handler"
	"music-metadata/internal/handlers/song_handler"
	"music-metadata/internal/middleware"
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/cover_service"
	"music-metadata/internal/service/genre_service"
	"music-metadata/internal/service/song_service"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))

	musicFilesClient := music_files_client.NewClient(ac.Config.HttpServer.MusicFilesAddress)
	audioFileClient := audio_file_client.NewAudioFileClient(musicFilesClient)

	albumRepo := album_repo.NewRepository()
	artistRepo := artist_repo.NewRepository()
	genreRepo := genre_repo.NewRepository()
	songRepo := song_repo.NewRepository()
	txManager := service.NewTransactionManager(*ac.Db)

	albumService := album_service.NewService(albumRepo)
	artistService := artist_service.NewService(artistRepo)
	genreService := genre_service.NewService(genreRepo)
	songService := song_service.NewService(songRepo, *albumService, *artistService, *genreService, audioFileClient)
	coverService := cover_service.NewService(*songService, audioFileClient)

	albumHandler := album_handler.NewHandler(*albumService, *coverService, txManager)
	artistHandler := artist_handler.NewHandler(*artistService, *coverService, txManager)
	genreHandler := genre_handler.NewHandler(*genreService, *coverService, txManager)
	songHandler := song_handler.NewHandler(*songService, txManager)

	api := r.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/scan", songHandler.Scan)

		album := api.Group("/albums")
		{
			album.GET("", albumHandler.GetAll)
		}

		artist := api.Group("/artists")
		{
			artist.GET("", artistHandler.GetAll)
		}

		genre := api.Group("/genres")
		{
			genre.GET("", genreHandler.GetAll)
		}
	}

	log.Debug().Msg("Router setup successfully")
	return r
}
