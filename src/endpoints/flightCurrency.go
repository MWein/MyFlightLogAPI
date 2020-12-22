package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type FlightCurrency struct {
	VFRDay               bool     `json:"vfrDay"`
	VFRDayTO             int      `json:"vfrDayTO"`
	VFRDayLandings       int      `json:"vfrDayLandings"`
	VFRNight             bool     `json:"vfrNight"`
	VFRNightTO           int      `json:"vfrNightTO"`
	VFRNightLandings     int      `json:"vfrNightLandings"`
	Instrument           bool     `json:"instrument"`
	InstrumentApproaches int      `json:"instrumentApproaches"`
	ApproachDates        []string `json:"lastApproaches"`
}

func FlightCurrencyRequirements(w http.ResponseWriter, r *http.Request) {
	var reqs FlightCurrency

	// 3 VFR takeoff and landings (Day or night)
	err := database.DBConnection.QueryRow("SELECT sum(takeoffs) AS takeoffs, sum(landings) AS landings, (sum(takeoffs) >= 3 AND sum(landings) >= 3) AS met FROM log WHERE NOW()::date - date <= 90").Scan(&reqs.VFRDayTO, &reqs.VFRDayLandings, &reqs.VFRDay)
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

	// Last 6 instrument approaches
	rows, _ := database.DBConnection.Query("SELECT date, instrument_approaches FROM log WHERE instrument_approaches > 0 ORDER BY date DESC LIMIT 6")
	for rows.Next() {
		var approachDate string
		var number int
		rows.Scan(&approachDate, &number)

		// Remove timestamp from date
		approachDate = approachDate[0:10]

		// Add copies of approach dates for each number returned by the query
		for i := 0; i < number; i++ {
			if len(reqs.ApproachDates) < 6 {
				reqs.ApproachDates = append(reqs.ApproachDates, approachDate)
			}
		}
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(reqs)
}
