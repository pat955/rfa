package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pat955/rss_feed_aggregator/api"
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
	r.HandleFunc("/v1/users", api.GetUser).Methods("GET")
	corsMux := middlewareLog(middlewareCors(r))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	srv.ListenAndServe()
}
