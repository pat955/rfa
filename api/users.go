package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

type User struct {
	Name string `json:"name"`
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var u User

	decodeForm(r, &u)
	apiConfig := connect()
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

func GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}

func GetApiKey(w http.ResponseWriter, r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		respondWithError(w, 403, "empty authorization header")
		return ""
	}

	apiString := strings.Split(auth, "ApiKey ")
	if len(apiString) != 2 || apiString[1] == "" {
		respondWithError(w, 403, "no apikey")
		return ""
	}
	return apiString[1]

}
