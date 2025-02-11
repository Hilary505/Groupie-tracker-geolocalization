package handlers

import (
	"html/template"
	"log"
	"net/http"

	"geolocalization/api"
)

var tmplt *template.Template

func init() {
	var err error
	tmplt, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Print("failed to parse templates: " + err.Error())
	}
}

// HomeHandler handles the home page requests.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if isValidPath(w, r) {
		// Handle valid paths here, e.g.:
		if r.URL.Path == "/" {
			// Serve home page
			http.ServeFile(w, r, "./templates/artist.html")
		}
		// else if r.URL.Path == "/artist" {
		// 	// Serve artist page
		// 	http.ServeFile(w, r, "./templates/artist.html")
		// }
	}
	artists, err := api.FetchArtists(w, r)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}

	// Execute the template with data
	if execErr := tmplt.ExecuteTemplate(w, "index.html", artists); execErr != nil {
		api.HandleError(w, execErr, http.StatusInternalServerError, "500.html")
		return
	}
}

// ArtistHandler handles requests for artist-specific information.
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.HandleError(w, nil, http.StatusMethodNotAllowed, "405.html")
		return
	}
	if !isValidPath(w, r) {
		err := http.ErrNotSupported // Define an appropriate error
		api.HandleError(w, err, http.StatusBadRequest, "404.html")
		return
	}

	artists, err := api.FetchArtists(w, r)
	if err != nil {
		return
	}

	// Execute the template with data
	if execErr := tmplt.ExecuteTemplate(w, "artist.html", artists); execErr != nil {
		api.HandleError(w, execErr, http.StatusInternalServerError, "500.html")
		return
	}
}

// isValidPath checks if the requested URL path is valid.
func isValidPath(w http.ResponseWriter, r *http.Request) bool {
	path := r.URL.Path // Get the path from the request
	if path != "/" && path != "/artist" {
		http.ServeFile(w, r, "./templates/404.html") // Serve the 404 page
		return false
	}
	return true
}
