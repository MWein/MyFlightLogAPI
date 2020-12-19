package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type InstrumentRequirements struct {
	CrossCountryPilotInCommandMet  bool    `json:"ccPicMet"`
	CrossCountryPilotInCommand     float64 `json:"ccPic"`
	InstrumentHoursMet             bool    `json:"instrHoursMet"`
	InstrumentHours                float64 `json:"instrHours"`
	InstrumentHoursWithCFIMet      bool    `json:"instrHoursWithCFIMet"`
	InstrumentHoursWithCFI         float64 `json:"instrHoursWithCFI"`
	LongCrossCountry               bool    `json:"longCC"`
	RecentInstrumentInstructionMet bool    `json:"recentInstrInstructionMet"`
	RecentInstrumentInstruction    float64 `json:"recentInstrInstruction"`
}

func InstrumentRatingProgress(w http.ResponseWriter, r *http.Request) {
	var reqs InstrumentRequirements

	// Cross country hours as PIC
	err := database.DBConnection.QueryRow("SELECT sum(total), sum(total) >= 50 AS met FROM log WHERE cross_country > 0 AND pic > 0").Scan(&reqs.CrossCountryPilotInCommand, &reqs.CrossCountryPilotInCommandMet)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Total instrument time
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument, (sum(sim_instrument) + sum(instrument)) >= 40 AS met FROM log").Scan(&reqs.InstrumentHours, &reqs.InstrumentHoursMet)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Total instrument time with CFI
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument, (sum(sim_instrument) + sum(instrument)) >= 40 AS met FROM log WHERE dual > 0").Scan(&reqs.InstrumentHoursWithCFI, &reqs.InstrumentHoursWithCFIMet)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Long cross country
	err = database.DBConnection.QueryRow("SELECT count(*) > 0 AS met FROM log WHERE remarks LIKE '%250nm Cross Country%'").Scan(&reqs.LongCrossCountry)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// 3 hours instrument within last 2 calendar months
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument, sum(sim_instrument) + sum(instrument) >= 3 AS met FROM log WHERE (sim_instrument > 0 OR instrument > 0) AND dual > 0 AND (extract(year from age(NOW()::date, date)) * 12 + extract(month from age(NOW()::date, date))) <= 2").Scan(&reqs.RecentInstrumentInstruction, &reqs.RecentInstrumentInstructionMet)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(reqs)
}
