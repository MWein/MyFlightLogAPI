package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type Airport struct {
	Ident string  `json:"ident"`
	Name  string  `json:"name"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	LastVisited string `json:"lastVisited"`
}

type Airports []Airport

func VisitedAirports(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DBConnection.Query("SELECT ident, name, lat, long, (SELECT date FROM log WHERE airport.ident = ANY(stops) ORDER BY date DESC LIMIT 1) FROM (SELECT unnest(stops) FROM log GROUP BY unnest) AS x JOIN airport ON ident = unnest")

	var airports Airports
	for rows.Next() {
		var airport Airport
		rows.Scan(&airport.Ident, &airport.Name, &airport.Lat, &airport.Long, &airport.LastVisited)

		// Remove timestamp from date
		airport.LastVisited = airport.LastVisited[0:10]

		airports = append(airports, airport)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(airports)
}
