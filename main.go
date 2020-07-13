package main

import (
	"log"
	"net/http"
	"github.com/MWein/MyFlightLogAPI/src/endpoints"
)


func handleRequests() {
	http.HandleFunc("/log", endpoints.FlightLogs)
	http.HandleFunc("/plane-image", endpoints.PlaneImg)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}