package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"geolocalization/models"
	// "tracker/models"
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
			lat, lon, geoErr := GeocodeLocation(loc.Locations)
			if geoErr != nil {
				fmt.Println("Geocoding error:", geoErr)
				continue
			}
			formatLocation := FormatLocations(loc.Locations)

			geocodedLocations = append(geocodedLocations, models.Location{
				Locations: formatLocation,
				Lat:       lat,
				Lon:       lon,
			})
			break // Stop looping after finding the first match
		}
	}
	wrapper = locationIndex.Index
	return wrapper, nil
}

func GeocodeLocation(locationNames []string) (lat, lon float64, err error) {
	baseURL := "https://geocode.search.hereapi.com/v1/geocode"
	apiKey := "8UIoooRk33BTptWdvimLIiIWA-Ss0T8LguDmhzb8-Xs"

	// Iterate through each location in the list
	for _, locationName := range locationNames {
		params := url.Values{}
		params.Set("q", locationName)
		params.Set("apiKey", apiKey)

		reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(reqURL)
		if err != nil {
			fmt.Printf("Error geocoding location '%s': %s\n", locationName, err) // Debug message for geocoding failure
			continue                                                             // If geocoding fails, continue with the next location
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
			continue // Skip to the next location if decoding fails
		}
		var geocodedLocations []string
		// If coordinates were found for the location, return them
		if len(result.Items) > 0 {
			lat  = result.Items[0].Position.Lat
			lon = result.Items[0].Position.Lon
			//fmt.Printf("Successfully geocoded '%s': %s, %s\n", locationName, lat, lon)
			geocodedLocations = append(geocodedLocations, fmt.Sprintf("%v,%v", lat, lon))
		} else {
			fmt.Printf("No coordinates found for '%s'.\n", locationName)
		}
	}

	
	return 0.0, 0.0 , nil
}

func FormatLocations(locations []string) []string {
	var formattedLocations []string
	for _, loc := range locations {
		formattedLocations = append(formattedLocations, fmt.Sprintf("%v", loc))
	}
	return formattedLocations
}
