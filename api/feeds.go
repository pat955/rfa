package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

type Feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

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
		respondWithError(w, 409, err.Error())
		return
	}
	respondWithJSON(w, 200, f)
}

func GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	allFeeds, err := Connect().DB.RetrieveFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, allFeeds)
}

func DeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var url URL
	decodeForm(r, &url)
	a := Connect()
	f, _ := a.DB.GetFeed(r.Context(), url.URL)
	fmt.Println(f.UserID, user.ID)
	if f.UserID == user.ID {
		a.DB.DeleteFeed(r.Context(), f.Url)
		respondWithJSON(w, 204, nil)
		return
	}
	respondWithError(w, 404, "feed not found")
}

type URL struct {
	URL string `json:"url"`
}
