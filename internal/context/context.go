package context

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"music-metadata/internal/config"
)

type AppContext struct {
	Config *config.Configuration
	Db     *sqlx.DB
	Logger *zerolog.Logger
}
