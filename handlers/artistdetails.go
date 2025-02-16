package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"geolocalization/api"
	"geolocalization/models"
)

// ArtistDetails handles requests for details about a specific artist.
func ArtistDetails(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method != http.MethodGet {
		api.HandleError(w, nil, http.StatusMethodNotAllowed, "405.html")
		return
	}

	artistIdStr := r.URL.Query().Get("artistId")
	artistId, err := strconv.Atoi(artistIdStr)
	if err != nil {
		api.HandleError(w, err, http.StatusBadRequest, "400.html")
		return
	}

	artists, err := api.FetchArtists(w, r)
	if err != nil {
		return
	}

	var selectedArtist models.Artist

	for _, artist := range artists {
		if artist.ID == artistId {
			selectedArtist = artist
			break
		}
	}

	if selectedArtist.ID == 0 {
		api.HandleError(w, err, http.StatusBadRequest, "400.html")
		return
	}

	// Fetch location data
	locations, err := api.FetchLocations(artistId, w, r)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}

	// Fetch concert dates
	concertDates, err := api.FetchConcertDates(w, r)
	if err != nil {
		fmt.Println("2")
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}

	// Fetch relations
	relations, err := api.FetchRelations(w, r)
	if err != nil {
		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}

	// Filter data by artist ID
	var artistLocations []models.Location
	for _, location := range locations {
		if location.ID == artistId {
			artistLocations = append(artistLocations, location)
		}
	}

	var artistConcertDates []models.Date
	for _, concertDate := range concertDates.Index {
		if concertDate.Id == artistId {
			artistConcertDates = append(artistConcertDates, concertDate)
		}
	}

	var artistRelations []models.Relation
	for _, relation := range relations.Index {
		if relation.ID == artistId {
			artistRelations = append(artistRelations, relation)
		}
	}

	// Combine artist data with related data
	data := struct {
		Artist       models.Artist
		Locations    []models.Location
		ConcertDates []models.Date
		Relations    []models.Relation
	}{
		Artist:       selectedArtist,
		Locations:    artistLocations,
		ConcertDates: artistConcertDates,
		Relations:    artistRelations,
	}

	// Execute the template with data
	err = tmplt.ExecuteTemplate(w, "artistdetails.html", data)
	if err != nil {

		api.HandleError(w, err, http.StatusInternalServerError, "500.html")
		return
	}
}
