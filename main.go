package main

import (
	"log"
	"net/http"

	"github.com/european-go/rating-system/handler"
)

const PORT = "9000"

func main() {

	router := http.NewServeMux()
	router.HandleFunc("POST /new_rating", handler.RouteNewRating)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}
	log.Println("Listening on port " + PORT + " ...")
	server.ListenAndServe()
}
