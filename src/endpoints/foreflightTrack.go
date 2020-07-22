package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

func ForeflightTrack(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["flightid"]

	if !ok || len(id[0]) < 1 {
		fmt.Fprintf(w, "Flight ID is required")
		return
	}

	rows, err := database.DBConnection.Query("SELECT lat, long FROM foreflight WHERE flightid = $1 ORDER BY sequence", id[0])
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	var foreflightTrack [][2]float64
	for rows.Next() {
		var latLon [2]float64
		rows.Scan(&latLon[0], &latLon[1])

		foreflightTrack = append(foreflightTrack, latLon)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(foreflightTrack)
}
