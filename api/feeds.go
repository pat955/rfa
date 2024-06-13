package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed Feed
	decodeForm(r, &feed)
	a := Connect()
	f, err := a.DB.CreateFeed(
		r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      feed.Name,
			Url:       feed.URL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, f)
}

type Feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
