package main

import (
	"log"
	"net/http"
	"os"

	"geolocalization/handlers"
)

func main() {
	if len(os.Args) > 1 {
		log.Println("The program epects only one argument: \n \n e.g go run main.go ")
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.ArtistHandler)
	mux.HandleFunc("/relations", handlers.RelationHandler)
	mux.HandleFunc("/locations", handlers.LocationHandler)
	mux.HandleFunc("/dates/", handlers.DatesHandler)
	mux.HandleFunc("/artistProfile", handlers.ArtistDetails)

	// Serve static files (CSS, images)
	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Start the server
	log.Println("Server running at  port: http://localhost:8000 \n ")
	log.Print("Click CTR + C  To terminate the server")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
