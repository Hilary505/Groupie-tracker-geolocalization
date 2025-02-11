package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"geolocalization/models"
)

// FetchRelations retrieves relations data from the API.
func FetchRelations(w http.ResponseWriter, r *http.Request) (relations models.RelationsResponse, e error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
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
		return models.RelationsResponse{}, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if jsonErr := json.Unmarshal(body, &relations); jsonErr != nil {
		return models.RelationsResponse{}, fmt.Errorf("failed to decode relations: %w", jsonErr)
	}

	return relations, nil
}
