package context

import (
	"github.com/jmoiron/sqlx"
	"music-metadata/internal/config"
)

type AppContext struct {
	Config *config.Configuration
	Db     *sqlx.DB
}
