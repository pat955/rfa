package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	r.HandleFunc("/v1/get_feed/{feedID}", api.GetFeed).Methods("GET")
	r.HandleFunc("/v1/posts", api.Auth(api.GetPostsByUser)).Methods("GET")

	r.HandleFunc("/v1/debug/feed_follows", api.GetAllFollowFeeds).Methods("GET")

	corsMux := logMW(corsMW(r))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           corsMux,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	go func() {
		for {
			// change name?
			api.RetrieveGroup()
			time.Sleep(time.Second * 60)
		}
	}()
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Support pagination of the endpoints that can return many items
// Support different options for sorting and filtering posts using query parameters
// Classify different types of feeds and posts (e.g. blog, podcast, video, etc.)
// Add a CLI client that uses the API to fetch and display posts, maybe it even allows you to read them in your terminal
// Scrape lists of feeds themselves from a third-party site that aggregates feed URLs
// Add support for other types of feeds (e.g. Atom, JSON, etc.)
// Add integration tests that use the API to create, read, update, and delete feeds and posts
// Add bookmarking or "liking" to posts
// Create a simple web UI that uses your backend API
