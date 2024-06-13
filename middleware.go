package main

import (
	"log"
	"net/http"

	"github.com/pat955/rss_feed_aggregator/api"
)

func corsMW(next http.Handler) http.Handler {
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

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authMW(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := api.GetApiKey(w, r)
		if apiKey == "" {
			return
		}
		apiConfig := api.Connect()
		u, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r, u)
	}
}
