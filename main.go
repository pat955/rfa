package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pat955/rss_feed_aggregator/api"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	public := "../public"
	port := os.Getenv("PORT")

	r := mux.NewRouter()
	defaultHandler := http.StripPrefix("/app", http.FileServer(http.Dir(public)))
	r.Handle("/app/*", defaultHandler)

	r.HandleFunc("/v1/healthz", api.Health).Methods("GET")
	r.HandleFunc("/v1/err", api.Error).Methods("GET")
	r.HandleFunc("/v1/users", api.AddUser).Methods("POST")
	r.HandleFunc("/v1/users", authMW(api.GetUser)).Methods("GET")
	corsMux := logMW(corsMW(r))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	srv.ListenAndServe()
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)
