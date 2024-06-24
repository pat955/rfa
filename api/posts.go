package api

import (
	"net/http"

	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func GetPostsByUser(w http.ResponseWriter, r *http.Request, u database.User) {
	db := connect().DB
	posts, err := db.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		u.ID,
		int32(5),
	})
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	respondWithJSON(w, 200, posts)
}
