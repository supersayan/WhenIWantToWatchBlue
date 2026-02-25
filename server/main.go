package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	server := NewServer()
	router := mux.NewRouter()

	router.Use(corsMiddleware)

	router.HandleFunc("/tv/{index:[0-9]+}", server.getTvHandler).Methods("GET")
	router.HandleFunc("/tv/{index:[0-9]+}/flip", server.toggleTvHandler).Methods("POST")
	router.HandleFunc("/tv/all", server.getAllTvHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
