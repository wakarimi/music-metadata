package api

import (
	"music-metadata/internal/client/music_files_client"
	"music-metadata/internal/client/music_files_client/audio_file_client"
	"music-metadata/internal/context"
	"music-metadata/internal/database/repository/album_repo"
	"music-metadata/internal/database/repository/artist_repo"
	"music-metadata/internal/database/repository/genre_repo"
	"music-metadata/internal/database/repository/song_repo"
	"music-metadata/internal/handlers/album_handler"
	"music-metadata/internal/handlers/artist_handler"
	"music-metadata/internal/handlers/cover_handler"
	"music-metadata/internal/handlers/genre_handler"
	"music-metadata/internal/handlers/song_handler"
	"music-metadata/internal/middleware"
	"music-metadata/internal/service"
	"music-metadata/internal/service/album_service"
	"music-metadata/internal/service/artist_service"
	"music-metadata/internal/service/cover_service"
	"music-metadata/internal/service/genre_service"
	"music-metadata/internal/service/song_service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))
	r.Use(middleware.CORSMiddleware())

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
	coverHandler := cover_handler.NewHandler(*coverService, txManager)

	api := r.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/scan", songHandler.Scan)

		songs := api.Group("/songs")
		{
			songs.GET("/:songId", songHandler.Get)
			songs.GET("", songHandler.GetAll)
		}

		album := api.Group("/albums")
		{
			album.GET("/:albumId", albumHandler.Get)
			album.GET("", albumHandler.GetAll)
			album.GET("/:albumId/songs", songHandler.GetByAlbumId)
			album.GET("/:albumId/covers", coverHandler.GetAllByAlbumId)
		}

		artist := api.Group("/artists")
		{
			artist.GET("/:artistId", artistHandler.Get)
			artist.GET("", artistHandler.GetAll)
			artist.GET("/:artistId/songs", songHandler.GetByArtistId)
		}

		genre := api.Group("/genres")
		{
			genre.GET("/:genreId", genreHandler.Get)
			genre.GET("", genreHandler.GetAll)
			genre.GET("/:genreId/songs", songHandler.GetByGenreId)
		}
	}

	log.Debug().Msg("Router setup successfully")
	return r
}
