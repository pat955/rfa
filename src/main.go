package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	public := "../public"
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)
	apiConfig := apiConfig{DB: dbQueries}
	fmt.Println(apiConfig)

	r := mux.NewRouter()
	defaultHandler := http.StripPrefix("/app", http.FileServer(http.Dir(public)))
	r.Handle("/app/*", defaultHandler)

	r.HandleFunc("/v1/healthz", handlerHealth).Methods("GET")
	r.HandleFunc("/v1/err", handleError).Methods("GET")
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	srv.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	DB *database.Queries
}
