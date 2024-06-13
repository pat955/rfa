package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, Status{Status: "ok"})
}

func Error(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")

}

type Status struct {
	Status string `json:"status"`
}
