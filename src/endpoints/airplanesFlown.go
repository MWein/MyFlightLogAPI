package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type Plane struct {
	Ident      string `json:"ident"`
	Type       string `json:"type"`
	Flights    int    `json:"flights"`
	LastFlight string `json:"lastFlight"`
}

type Planes []Plane

func AirplanesFlown(w http.ResponseWriter, r *http.Request) {
	planeQuery := `SELECT
		ident,
		plane_type.name AS type,
		(SELECT count(*) FROM log WHERE log.ident = plane.ident) AS flights,
		(SELECT date FROM log WHERE log.ident = plane.ident ORDER BY date DESC LIMIT 1) AS last_flight
	FROM plane
	JOIN plane_type USING (type_id)`

	rows, _ := database.DBConnection.Query(planeQuery)

	var planes Planes
	for rows.Next() {
		var plane Plane
		rows.Scan(&plane.Ident, &plane.Type, &plane.Flights, &plane.LastFlight)

		// Remove timestamp from date
		plane.LastFlight = plane.LastFlight[0:10]

		planes = append(planes, plane)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(planes)
}
