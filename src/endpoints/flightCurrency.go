package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type FlightCurrency struct {
	VFRDay               bool `json:"vfrDay"`
	VFRDayTO             int  `json:"vfrDayTO"`
	VFRDayLandings       int  `json:"vfrDayLandings"`
	VFRNight             bool `json:"vfrNight"`
	VFRNightTO           int  `json:"vfrNightTO"`
	VFRNightLandings     int  `json:"vfrNightLandings"`
	Instrument           bool `json:"instrument"`
	InstrumentApproaches int  `json:"instrumentApproaches"`
}

func FlightCurrencyRequirements(w http.ResponseWriter, r *http.Request) {
	var reqs FlightCurrency

	// 3 VFR day takeoff and landings
	err := database.DBConnection.QueryRow("SELECT sum(takeoffs) AS takeoffs, sum(landings) AS landings, (sum(takeoffs) >= 3 AND sum(landings) >= 3) AS met FROM log WHERE NOW()::date - date <= 90 AND night = 0").Scan(&reqs.VFRDayTO, &reqs.VFRDayLandings, &reqs.VFRDay)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// 3 VFR night takeoff and landings
	err = database.DBConnection.QueryRow("SELECT sum(takeoffs) AS takeoffs, sum(landings) AS landings, (sum(takeoffs) >= 3 AND sum(landings) >= 3) AS met FROM log WHERE NOW()::date - date <= 90 AND night > 0").Scan(&reqs.VFRNightTO, &reqs.VFRNightLandings, &reqs.VFRNight)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// 6 Instrument Approaches
	err = database.DBConnection.QueryRow("SELECT sum(instrument_approaches), sum(instrument_approaches) >= 6 AS met FROM log WHERE (extract(year from age(NOW()::date, date)) * 12 + extract(month from age(NOW()::date, date))) <= 5").Scan(&reqs.InstrumentApproaches, &reqs.Instrument)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(reqs)
}
