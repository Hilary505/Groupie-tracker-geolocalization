package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"geolocalization/models"
	// "tracker/models"
)

// FetchArtists retrieves the list of artists from the API.
func FetchArtists(w http.ResponseWriter, r *http.Request) (artists []models.Artist, e error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		HandleError(w, err, 500, "500.html")
		e = err
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Println("Error closing response body:", closeErr)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if jsonErr := json.Unmarshal(body, &artists); jsonErr != nil {
		return nil, fmt.Errorf("failed to decode artists: %w", jsonErr)
	}
	return artists, nil
}
