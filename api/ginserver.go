package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GinServer struct {
	engine *gin.Engine
}

func (g *GinServer) Run(address string) (err error) {
	err = g.engine.Run(address)
	if err != nil {
		log.Error().Err(err).Msg("Failed to run Gin server")
	}
	return err
}
