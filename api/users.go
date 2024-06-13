package api

import (
	"errors"
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := getApiKey(r)
	if err != nil {
		respondWithError(w, 403, err.Error())
		return
	}
	apiConfig := Connect()
	user, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	respondWithJSON(w, 200, user)

}

func getApiKey(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	apiString := strings.Split(auth, "ApiKey ")
	if len(apiString) != 2 || apiString[1] == "" {
		return "", errors.New("authorization invalid")
	}
	return apiString[1], nil

}
