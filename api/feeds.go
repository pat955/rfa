package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

type FeedForJSON struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

func dbFeedToFeed(feed database.Feed) FeedForJSON {
	return FeedForJSON{ID: feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
}
func dbFeedsToFeeds(feeds []database.Feed) []FeedForJSON {
	updatedFeeds := make([]FeedForJSON, 0)
	for _, feed := range feeds {
		updatedFeeds = append(updatedFeeds, dbFeedToFeed(feed))
	}
	return updatedFeeds
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
	response := CreateFeedResponse{Feed: dbFeedToFeed(f), FeedFollow: follow_feed}
	respondWithJSON(w, 200, response)
}

func GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	allFeeds, err := connect().DB.RetrieveFeeds(r.Context())
	if err != nil {
		respondWithError(w, 404, "No feeds currently")
		return
	}
	respondWithJSON(w, 200, dbFeedsToFeeds(allFeeds))
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

func GetFeed(w http.ResponseWriter, r *http.Request) {
	feed_id, found := mux.Vars(r)["feedID"]
	if !found {
		respondWithError(w, 400, "no feed follow id in url")
		return
	}
	db := connect().DB
	feed, err := db.GetFeed(r.Context(), uuid.MustParse(feed_id))
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	newTime := time.Now().UTC()
	feed.UpdatedAt = newTime
	feed.LastFetchedAt = sql.NullTime{Time: newTime, Valid: true}
	err = db.MarkedFetched(r.Context(), database.MarkedFetchedParams{
		ID:            feed.ID,
		UpdatedAt:     newTime,
		LastFetchedAt: feed.LastFetchedAt,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, dbFeedToFeed(feed))
}
