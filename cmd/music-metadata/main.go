package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"music-metadata/api"
	"music-metadata/internal/config"
	"music-metadata/internal/database"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger().
		With().Str("service", "music-files").Logger()

	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
		return
	}

	log.Logger = log.Level(zerolog.DebugLevel)

	db, err := database.ConnectDb(cfg.DatabaseConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
	database.SetDatabase(db)

	if err = database.RunMigrations(db, "./internal/database/migrations"); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
		return
	}

	r := api.SetupRouter(cfg, db)
	if err = r.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
