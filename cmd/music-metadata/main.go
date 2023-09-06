package main

import (
	"log"
	"music-metadata/api"
	"music-metadata/internal/config"
	"music-metadata/internal/database"
)

func main() {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	log.Println(cfg.DatabaseConnectionString)
	db, err := database.ConnectDb(cfg.DatabaseConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	database.SetDatabase(db)

	err = database.RunMigrations(db, "./internal/database/migrations")
	if err != nil {
		log.Fatal("Failed to run migrations: %v", err)
	}

	r := api.SetupRouter()
	r.Run(":" + cfg.Port)
}
