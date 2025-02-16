package models

// Define the data structures for the API responses
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int       `json:"id"`
	Locations []string  `json:"locations"`
	Date      string    `json:"dates"`
	Lon       []float64 `json:"lat"`
	Lat       []float64 `json:"lon"`
}

type LocationsResponse struct {
	Index []Location `json:"index"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesResponse struct {
	Index []Date `json:"index"`
}

type Relation struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}
