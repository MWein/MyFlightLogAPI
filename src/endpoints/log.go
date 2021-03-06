package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
	"github.com/lib/pq"
)

type Log struct {
	Id                   string       `json:"id"`
	Date                 string       `json:"date"`
	Type                 string       `json:"type"`
	Ident                string       `json:"ident"`
	Stops                []string     `json:"stops"`
	InstrumentApproaches int          `json:"instrumentApproaches"`
	Night                float32      `json:"night"`
	Instrument           float32      `json:"instrument"`
	Sim_instrument       float32      `json:"simInstrument"`
	Flight_sim           float32      `json:"flightSim"`
	Cross_country        float32      `json:"crossCountry"`
	Instructor           float32      `json:"instructor"`
	Dual                 float32      `json:"dual"`
	Pic                  float32      `json:"pic"`
	Total                float32      `json:"total"`
	Takeoffs             int          `json:"takeoffs"`
	Landings             int          `json:"landings"`
	Remarks              string       `json:"remarks"`
	Geolocation          [][2]float64 `json:"geolocation"`
	Pictures             []string     `json:"pictures"`
	HasFFTrack           bool         `json:"hasForeflightTrack"`
	Favorite             bool         `json:"favorite"`
}
type Logs []Log

func arrayIncludes(a string, list []string) bool {
	for _, item := range list {
		if a == item {
			return true
		}
	}

	return false
}

func FlightLogs(w http.ResponseWriter, r *http.Request) {
	// Retrieve main logs body

	logsStatement := `SELECT
	id, date, plane_type.name AS type, ident, stops, instrument_approaches, night, instrument, sim_instrument, flight_sim, cross_country, instructor, dual, log.pic, total, takeoffs, landings, remarks, favorite
	FROM log
	JOIN plane USING (ident)
	JOIN plane_type USING (type_id)
	ORDER BY date DESC`

	rows, _ := database.DBConnection.Query(logsStatement)

	var logs Logs
	for rows.Next() {
		var log Log
		rows.Scan(&log.Id, &log.Date, &log.Type, &log.Ident, pq.Array(&log.Stops), &log.InstrumentApproaches, &log.Night, &log.Instrument, &log.Sim_instrument, &log.Flight_sim, &log.Cross_country, &log.Instructor, &log.Dual, &log.Pic, &log.Total, &log.Takeoffs, &log.Landings, &log.Remarks, &log.Favorite)

		// Remove timestamp from date
		log.Date = log.Date[0:10]

		logs = append(logs, log)
	}

	// Extract all unique airport codes
	var airportCodes []string
	for _, log := range logs {
		for _, code := range log.Stops {
			if !arrayIncludes(code, airportCodes) {
				airportCodes = append(airportCodes, code)
			}
		}
	}

	// Generate airport map
	geoLocationsStatement := `SELECT ident, lat, long FROM airport WHERE ident = ANY($1)`
	rows, _ = database.DBConnection.Query(geoLocationsStatement, pq.Array(airportCodes))

	airportMap := make(map[string][2]float64)

	for rows.Next() {
		var ident string
		var latLon [2]float64
		rows.Scan(&ident, &latLon[0], &latLon[1])
		airportMap[ident] = latLon
	}

	for i := 0; i < len(logs); i++ {
		log := &logs[i]
		stops := log.Stops
		log.Pictures = []string{}

		// Generate geolocations and add to logs
		for _, stop := range stops {
			log.Geolocation = append(log.Geolocation, airportMap[stop])
		}

		// Add picture IDs
		rows, err := database.DBConnection.Query("SELECT id FROM pictures WHERE flightid = $1", log.Id)
		if err != nil {
			continue
		}

		for rows.Next() {
			var pic string
			rows.Scan(&pic)
			log.Pictures = append(log.Pictures, pic)
		}

		// Add boolean for Foreflight Track
		err = database.DBConnection.QueryRow("SELECT count(*) > 0 FROM foreflight WHERE flightid = $1", log.Id).Scan(&log.HasFFTrack)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(logs)
}
