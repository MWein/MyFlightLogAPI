package main

import (
	"log"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
	"github.com/MWein/MyFlightLogAPI/src/endpoints"
)

func handleRequests() {
	http.HandleFunc("/log", endpoints.FlightLogs)
	http.HandleFunc("/plane-image", endpoints.PlaneImg)
	http.HandleFunc("/flight-image", endpoints.FlightImg)
	http.HandleFunc("/foreflight-track", endpoints.ForeflightTrack)
	http.HandleFunc("/visited-airports", endpoints.VisitedAirports)
	http.HandleFunc("/airplanes-flown", endpoints.AirplanesFlown)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	database.StartDB()

	handleRequests()
}
