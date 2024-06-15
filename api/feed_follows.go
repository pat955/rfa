package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

type FeedFollow struct {
	FeedID string `json:"feed_id"`
}

func GetAllFollowFeeds(w http.ResponseWriter, r *http.Request) {
	a := connect()
	if feed_follows, err := a.DB.GetAllFeedFollows(r.Context()); err != nil {
		respondWithError(w, 404, "no feed follows found")
	} else {
		respondWithJSON(w, 200, feed_follows)
	}
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
		respondWithError(w, 409, "already following this feed")
		return
	}
	respondWithJSON(w, 200, newFeedFollow)
}

// feed_id != feed_follow_id
func UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follow_id, found := mux.Vars(r)["feedFollowID"]
	if !found {
		respondWithError(w, 400, "no feed follow id in url")
		return
	}
	id, err := uuid.Parse(feed_follow_id)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	a := connect()
	err = a.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "no feed to unfollow")
		return
	}
	respondWithJSON(w, 204, nil)
}

func GetAllFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	a := connect()
	allFollowed, err := a.DB.GetAllFollowed(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, "no feeds followed by user")
		return
	}
	respondWithJSON(w, 200, allFollowed)
}
