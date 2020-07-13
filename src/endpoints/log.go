package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

type Log struct {
	Id             string       `json:"id"`
	Date           string       `json:"date"`
	Type           string       `json:"type"`
	Ident          string       `json:"ident"`
	Stops          []string     `json:"stops"`
	Night          float32      `json:"night"`
	Instrument     float32      `json:"instrument"`
	Sim_instrument float32      `json:"simInstrument"`
	Flight_sim     float32      `json:"flightSim"`
	Cross_country  float32      `json:"crossCountry"`
	Instructor     float32      `json:"instructor"`
	Dual           float32      `json:"dual"`
	Pic            float32      `json:"pic"`
	Total          float32      `json:"total"`
	Takeoffs       int          `json:"takeoffs"`
	Landings       int          `json:"landings"`
	Remarks        string       `json:"remarks"`
	Geolocation    [][2]float64 `json:"geolocation"`
	Pictures       []string     `json:"pictures"`
	HasFFTrack     bool         `json:"hasForeflightTrack"`
}
type Logs []Log

type Totals struct {
	Takeoffs      int     `json:"takeoffs"`
	Landings      int     `json:"landings"`
	Night         float32 `json:"night"`
	Instrument    float32 `json:"instrument"`
	SimInstrument float32 `json:"simInstrument"`
	FlightSim     float32 `json:"flightSim"`
	CrossCountry  float32 `json:"crossCountry"`
	Instructor    float32 `json:"instructor"`
	Dual          float32 `json:"dual"`
	Pic           float32 `json:"pic"`
	Total         float32 `json:"total"`
}

type LogsReturn struct {
	Logs   Logs   `json:"logs"`
	Totals Totals `json:"totals"`
}

func arrayIncludes(a string, list []string) bool {
	for _, item := range list {
		if a == item {
			return true
		}
	}

	return false
}

func FlightLogs(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Retrieve main logs body

	logsStatement := `SELECT
	id, date, plane_type.name AS type, ident, stops, night, instrument, sim_instrument, flight_sim, cross_country, instructor, dual, log.pic, total, takeoffs, landings, remarks
	FROM log
	JOIN plane USING (ident)
	JOIN plane_type USING (type_id)`

	rows, err := db.Query(logsStatement)

	var logs Logs
	for rows.Next() {
		var log Log
		rows.Scan(&log.Id, &log.Date, &log.Type, &log.Ident, pq.Array(&log.Stops), &log.Night, &log.Instrument, &log.Sim_instrument, &log.Flight_sim, &log.Cross_country, &log.Instructor, &log.Dual, &log.Pic, &log.Total, &log.Takeoffs, &log.Landings, &log.Remarks)
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
	rows, err = db.Query(geoLocationsStatement, pq.Array(airportCodes))

	airportMap := make(map[string][2]float64)

	for rows.Next() {
		var ident string
		var latLon [2]float64
		rows.Scan(&ident, &latLon[0], &latLon[1])
		airportMap[ident] = latLon
	}

	// Generate geolocations and add to logs
	for i := 0; i < len(logs); i++ {
		log := &logs[i]
		stops := log.Stops

		for _, stop := range stops {
			log.Geolocation = append(log.Geolocation, airportMap[stop])
		}
	}



	// TODO add pictures
	for i := 0; i < len(logs); i++ {
		log := &logs[i]
		
		log.Pictures = []string{}
	}



	// Retrieve totals

	totalsStatement := `SELECT sum(takeoffs) AS takeoffs, sum(landings) AS landings, sum(night) AS night, sum(instrument) AS instrument, sum(sim_instrument) AS sim_instrument, sum(flight_sim) AS flight_sim, sum(cross_country) AS cross_country, sum(instructor) AS instructor, sum(dual) AS dual, sum(pic) AS pic, sum(total) AS total FROM log`

	var totals Totals
	err = db.QueryRow(totalsStatement).Scan(&totals.Takeoffs, &totals.Landings, &totals.Night, &totals.Instrument, &totals.SimInstrument, &totals.FlightSim, &totals.CrossCountry, &totals.Instructor, &totals.Dual, &totals.Pic, &totals.Total)

	returnValue := LogsReturn{
		Logs:   logs,
		Totals: totals,
	}

	json.NewEncoder(w).Encode(returnValue)
}
