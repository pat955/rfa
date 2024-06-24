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
	r.HandleFunc("/v1/users", api.Auth(api.GetUser)).Methods("GET")
	r.HandleFunc("/v1/feeds", api.Auth(api.CreateFeed)).Methods("POST")
	r.HandleFunc("/v1/feeds", api.GetAllFeeds).Methods("GET")
	r.HandleFunc("/v1/feeds", api.Auth(api.DeleteFeed)).Methods("DELETE")
	r.HandleFunc("/v1/feed_follows", api.Auth(api.FollowFeed)).Methods("POST")
	r.HandleFunc("/v1/feed_follows", api.Auth(api.GetAllFollowedFeeds)).Methods("POST")
	r.HandleFunc("/v1/feed_follows/{feedFollowID}", api.Auth(api.UnfollowFeed)).Methods("DELETE")
	r.HandleFunc("/v1/next_to_fetch", api.GetNextFeedsToFetch).Methods("GET")

	r.HandleFunc("/v1/debug/feed_follows", api.GetAllFollowFeeds).Methods("GET")

	corsMux := logMW(corsMW(r))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	srv.ListenAndServe()
}
