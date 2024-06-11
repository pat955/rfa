package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("..env")
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.HandleFunc("v1/healtz", handlerHealth).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	srv.ListenAndServe()
}
