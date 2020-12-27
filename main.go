package main

import (
	"log"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
	"github.com/MWein/MyFlightLogAPI/src/endpoints"
)

func handleRequests() {
	// Logbook
	http.HandleFunc("/log", endpoints.FlightLogs)
	http.HandleFunc("/foreflight-track", endpoints.ForeflightTrack)

	// Upload Foreflight Track
	http.HandleFunc("/add-foreflight-track", endpoints.AddForeflightTrack)

	// Airplanes
	http.HandleFunc("/airplanes-flown", endpoints.AirplanesFlown)

	// Airports
	http.HandleFunc("/visited-airports", endpoints.VisitedAirports)

	// Build Log
	http.HandleFunc("/build-projects", endpoints.BuildProjects)
	http.HandleFunc("/build-details", endpoints.BuildDetails)

	// Currency Endpoints
	http.HandleFunc("/flight-currency", endpoints.FlightCurrencyRequirements)
	http.HandleFunc("/instrument-rating-progress", endpoints.InstrumentRatingProgress)

	// Consolidated Image Endpoint
	http.HandleFunc("/img", endpoints.Img)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	database.StartDB()

	handleRequests()
}
