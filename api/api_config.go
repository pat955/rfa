package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

var dbURL string

type apiConfig struct {
	DB *database.Queries
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	dbURL = os.Getenv("CONNECTION_STRING")
}

// dburl
func Connect() *apiConfig {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)
	return &apiConfig{DB: dbQueries}
}
