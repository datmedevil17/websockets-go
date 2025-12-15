package main

import (
	"log"
	"net/http"

	"collab-doc/server"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	hub := server.NewHub()
	go hub.Run()

	server.RegisterRoutes(r, hub)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
