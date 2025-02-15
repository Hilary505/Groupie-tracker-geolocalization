package handlers

import (
	"net/http"
	"strconv"

	"geolocalization/api"
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

	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}

	if err := tmplt.ExecuteTemplate(w, "locations.html", locations); err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
	// Set content type before writing response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
