package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"geolocalization/models"
)

// FetchLocations retrieves location data from the API.
func FetchLocations(artistID int, w http.ResponseWriter, r *http.Request) (wrapper []models.Location, e error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
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

	var locationIndex models.LocationsResponse
	if jsonErr := json.Unmarshal(body, &locationIndex); jsonErr != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", jsonErr)
	}
	// Find the artist's locations
	var geocodedLocations []models.Location
	for _, loc := range locationIndex.Index {
		if loc.ID == artistID { // Only process the requested artist
			allLat, allLon, geoErr := GeocodeLocation(loc.Locations)
			if geoErr != nil {
				fmt.Println("Geocoding error:", geoErr)
				continue
			}
			formatLocation := FormatLocations(loc.Locations)

			geocodedLocations = append(geocodedLocations, models.Location{
				Locations: formatLocation,
				Lat:       allLat,
				Lon:       allLon,
			})
		}
	}
	wrapper = geocodedLocations
	return wrapper, nil
}

func GeocodeLocation(locationNames []string) (alllat, alllon []float64, err error) {
	baseURL := "https://geocode.search.hereapi.com/v1/geocode"
	apiKey := "8UIoooRk33BTptWdvimLIiIWA-Ss0T8LguDmhzb8-Xs" // Replace with your actual API key

	var allLat []float64
	var allLon []float64

	for _, locationName := range locationNames {
		params := url.Values{}
		params.Set("q", locationName)
		params.Set("apiKey", apiKey)

		reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(reqURL)
		if err != nil {
			fmt.Printf("Error geocoding location '%s': %s\n", locationName, err)
			continue
		}
		defer resp.Body.Close()

		var result struct {
			Items []struct {
				Position struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lng"`
				} `json:"position"`
			} `json:"items"`
		}

		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("Error decoding geocoding response for '%s': %s\n", locationName, err)
			continue
		}

		if len(result.Items) > 0 {
			lat := result.Items[0].Position.Lat
			lon := result.Items[0].Position.Lon
			allLat = append(allLat, lat)
			allLon = append(allLon, lon)
		} else {
			fmt.Printf("No coordinates found for '%s'\n", locationName)
		}
	}

	return allLat, allLon, nil

}

func FormatLocations(locations []string) []string {
	var formattedLocations []string
	for _, loc := range locations {
		formattedLocations = append(formattedLocations, fmt.Sprintf("%v", loc))
	}
	return formattedLocations
}
