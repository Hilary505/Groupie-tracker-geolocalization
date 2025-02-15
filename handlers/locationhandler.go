package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"geolocalization/api"
	"geolocalization/models"
)

// LocationHandler handles requests for locations associated with a specific artist.
func LocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.HandleError(w, nil, http.StatusMethodNotAllowed, "405.html")
		return
	}

	artistIdStr := r.URL.Query().Get("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil || artistId > 56 {
		api.HandleError(w, err, http.StatusBadRequest, "400.html")
		return
	}

	locations, err := api.FetchLocations(artistId, w, r)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
	fmt.Println("Locations before filtering:", locations)

	filteredLocations := filterLocationsByArtistID(locations, artistId)
	fmt.Println("Filtered Locations:", filteredLocations)
	locationsJSON, err := json.Marshal(filteredLocations)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
	fmt.Println("Locations JSON:", string(locationsJSON))


	if err := tmplt.ExecuteTemplate(w, "locations.html", string(locationsJSON)); err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
	// Set content type before writing response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// filterLocationsByArtistID filters locations based on the artist ID.
func filterLocationsByArtistID(locations []models.Location, artistId int) []models.Location {
	var filtered []models.Location
	for _, location := range locations {
		if location.ID == artistId {
			filtered = append(filtered, location)
		}
	}
	return filtered
}
