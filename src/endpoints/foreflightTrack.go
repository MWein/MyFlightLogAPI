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

	var foreflightTrackBytes []byte
	database.DBConnection.QueryRow("SELECT data FROM foreflight WHERE flightid = $1", id[0]).Scan(&foreflightTrackBytes)

	var foreflightTrack [][2]float64
	json.Unmarshal(foreflightTrackBytes, &foreflightTrack)

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(foreflightTrack)
}
