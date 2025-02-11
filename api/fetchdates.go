package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"geolocalization/models"
	// "tracker/models"
)

// FetchConcertDates retrieves concert dates from the API.
func FetchConcertDates(w http.ResponseWriter, r *http.Request) (concertDates models.DatesResponse, e error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
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
		return models.DatesResponse{}, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if jsonErr := json.Unmarshal(body, &concertDates); jsonErr != nil {
		return models.DatesResponse{}, fmt.Errorf("failed to decode concert dates: %w", jsonErr)
	}
	return concertDates, nil
}
