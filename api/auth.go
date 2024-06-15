package api

import (
	"net/http"

	"github.com/pat955/rss_feed_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func Auth(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := GetApiKey(w, r)
		if apiKey == "" {
			return
		}
		apiConfig := connect()
		u, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r, u)
	}
}
