package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	var u User

	decodeForm(r, &u)
	apiConfig := Connect()
	newUser, err := apiConfig.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      u.Name,
		})
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	respondWithJSON(w, 200, newUser)
}

type User struct {
	Name string `json:"name"`
}
