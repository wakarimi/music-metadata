package api

import (
	"github.com/gin-gonic/gin"
	"music-metadata/internal/config"
	"music-metadata/internal/handlers"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/music-metadata-service")
	{
		api.POST("/scan", func(c *gin.Context) { handlers.Scan(c, cfg) })
	}

	return r
}
