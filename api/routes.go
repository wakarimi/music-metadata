package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/config"
	"music-metadata/internal/database/repository"
	"music-metadata/internal/handlers"
)

func SetupRouter(cfg *config.Configuration, db *sqlx.DB) *gin.Engine {
	albumRepo := repository.NewAlbumRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	trackMetadataRepo := repository.NewTrackMetadataRepository(db)
	musicHandler := handlers.NewMusicHandler(albumRepo, artistRepo, genreRepo, trackMetadataRepo)

	r := gin.Default()

	api := r.Group("/api/music-metadata-service")
	{
		api.POST("/scan", func(c *gin.Context) { musicHandler.Scan(c, cfg) })
	}

	return r
}
