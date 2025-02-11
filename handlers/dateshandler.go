package handlers

import (
	"net/http"
	"strconv"

	"geolocalization/api"
	"geolocalization/models"
)

// ConcertDatesHandler handles requests for concert dates associated with a specific artist.
func DatesHandler(w http.ResponseWriter, r *http.Request) {
	artistIdStr := r.URL.Query().Get("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil || artistId > 56 {
		api.HandleError(w, err, http.StatusBadRequest, "400.html")
		return
	}
	Dates, err := api.FetchConcertDates(w, r)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
	filteredDates := filterConcertDatesByArtistID(Dates.Index, artistId)
	if err := tmplt.ExecuteTemplate(w, "dates.html", filteredDates); err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
}

// filterConcertDatesByArtistID filters concert dates based on the artist ID.
func filterConcertDatesByArtistID(dates []models.Date, artistId int) []models.Date {
	var filtered []models.Date
	for _, date := range dates {
		if date.Id == artistId {
			filtered = append(filtered, date)
		}
	}
	return filtered
}
