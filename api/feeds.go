package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

type Feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FeedFollow struct {
	FeedID string `json:"feed_id"`
}
type ID struct {
	ID string `json:"id"`
}

func CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed Feed
	decodeForm(r, &feed)
	a := connect()
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
	allFeeds, err := connect().DB.RetrieveFeeds(r.Context())
	if err != nil {
		respondWithError(w, 404, "No feeds currently")
		return
	}
	respondWithJSON(w, 200, allFeeds)
}

func DeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed_id ID
	decodeForm(r, &feed_id)
	a := connect()

	id, err := uuid.Parse(feed_id.ID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	f, err := a.DB.GetFeed(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	if f.UserID == user.ID {
		a.DB.DeleteFeed(r.Context(), id)
		respondWithJSON(w, 204, nil)
		return
	}
	respondWithError(w, 403, "no permission")
}

func FollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed_follow FeedFollow
	decodeForm(r, &feed_follow)
	id, err := uuid.Parse(feed_follow.FeedID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	a := connect()
	feed, err := a.DB.GetFeed(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	newFeedFollow, err := a.DB.AddFeedFollow(
		r.Context(),
		database.AddFeedFollowParams{
			ID:        uuid.New(),
			FeedID:    feed.ID,
			UserID:    user.ID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, newFeedFollow)
}
