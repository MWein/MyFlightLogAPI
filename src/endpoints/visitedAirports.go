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
}

type Airports []Airport

func VisitedAirports(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DBConnection.Query("SELECT ident, name, lat, long FROM (SELECT unnest(stops) FROM log GROUP BY unnest) AS x JOIN airport ON ident = unnest")

	var airports Airports
	for rows.Next() {
		var airport Airport
		rows.Scan(&airport.Ident, &airport.Name, &airport.Lat, &airport.Long)
		airports = append(airports, airport)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(airports)
}
