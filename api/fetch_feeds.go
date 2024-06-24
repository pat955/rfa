package api

import (
	"net/http"
)

type Amount struct {
	Amount int32 `json:"amount"`
}

func GetNextFeedsToFetch(w http.ResponseWriter, r *http.Request) {
	var n Amount
	decodeForm(r, &n)
	db := connect().DB

	feedsToUpdate, err := db.GetNextFeedsToFetch(r.Context(), n.Amount)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, dbFeedsToFeeds(feedsToUpdate))
}
