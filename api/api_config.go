package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

// since i cannot load .env vars as a constant i have made it a psuedo const
var DBURL string

// Entry point to the DB, should maybe move it to internal for security?
type apiConfig struct {
	DB *database.Queries
}

// Loads DBURL, connection string
func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	DBURL = os.Getenv("CONNECTION_STRING")
}

// use this to connect to the db
func connect() *apiConfig {
	db, err := sql.Open("postgres", DBURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)
	return &apiConfig{DB: dbQueries}
}
