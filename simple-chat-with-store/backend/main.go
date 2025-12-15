package main

import (
	"log"
	"net/http"

	"github.com/datmedevil17/gorilla_sockets/server"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	h := server.NewHub()
	go h.Run()

	server.RegisterRoutes(r, h)

	address := ":8080"
	log.Println("Server listening on", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(err)
	}

}
