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

type ID struct {
	ID string `json:"id"`
}

type CreateFeedResponse struct {
	Feed       FeedForJSON         `json:"feed"`
	FeedFollow database.FeedFollow `json:"feed_follow"`
}
type Amount struct {
	Amount int32 `json:"amount"`
}

type FeedForJSON struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

func GetNextFeedsToFetch(w http.ResponseWriter, r *http.Request) {
	var n Amount
	decodeForm(r, &n)
	db := connect().DB
	feedsToUpdate, err := db.GetNextFeedsToFetch(r.Context(), n.Amount)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}
	respondWithJSON(w, 200, feedsToUpdate)
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
	// default follows feed if user created it
	follow_feed, err := a.DB.AddFeedFollow(r.Context(), database.AddFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    f.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, "error in createFeed() "+err.Error())
	}
	response := CreateFeedResponse{Feed: databaseFeedToFeed(f), FeedFollow: follow_feed}
	respondWithJSON(w, 200, response)
}

func GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	allFeeds, err := connect().DB.RetrieveFeeds(r.Context())
	if err != nil {
		respondWithError(w, 404, "No feeds currently")
		return
	}
	allFeedsResponse := make([]FeedForJSON, 0)
	for _, feed := range allFeeds {
		allFeedsResponse = append(allFeedsResponse, databaseFeedToFeed(feed))
	}
	respondWithJSON(w, 200, allFeedsResponse)
}

func databaseFeedToFeed(feed database.Feed) FeedForJSON {
	return FeedForJSON{ID: feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
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
