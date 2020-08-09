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
	http.HandleFunc("/flight-image", endpoints.FlightImg)
	http.HandleFunc("/foreflight-track", endpoints.ForeflightTrack)

	// Airplanes
	http.HandleFunc("/plane-image", endpoints.PlaneImg)
	http.HandleFunc("/airplanes-flown", endpoints.AirplanesFlown)

	// Airports
	http.HandleFunc("/visited-airports", endpoints.VisitedAirports)

	// Build Log
	http.HandleFunc("/build-projects", endpoints.BuildProjects)
	http.HandleFunc("/build-cover", endpoints.BuildCoverPhoto)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	database.StartDB()

	handleRequests()
}
