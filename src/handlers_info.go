package main

import "net/http"

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, Status{Status: "ok"})
}

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")

}

type Status struct {
	Status string `json:"status"`
}
