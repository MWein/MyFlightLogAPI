package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type InstrumentRequirements struct {
	InstrumentKnowledgeTest     bool    `json:"knowledgeTest"`
	CrossCountryPilotInCommand  float64 `json:"ccPic"`
	InstrumentHours             float64 `json:"instrHours"`
	InstrumentHoursWithCFI      float64 `json:"instrHoursWithCFI"`
	LongCrossCountry            bool    `json:"longCC"`
	RecentInstrumentInstruction float64 `json:"recentInstrInstruction"`
}

func InstrumentRatingProgress(w http.ResponseWriter, r *http.Request) {
	var reqs InstrumentRequirements

	// TODO Uncomment once I pass the knowledge test. No sense in making a DB change just for this.
	//reqs.InstrumentKnowledgeTest = true

	// Cross country hours as PIC
	err := database.DBConnection.QueryRow("SELECT sum(total) FROM log WHERE cross_country > 0 AND pic > 0").Scan(&reqs.CrossCountryPilotInCommand)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Total instrument time
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument FROM log").Scan(&reqs.InstrumentHours)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Total instrument time with CFI
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument FROM log WHERE dual > 0").Scan(&reqs.InstrumentHoursWithCFI)
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
	err = database.DBConnection.QueryRow("SELECT sum(sim_instrument) + sum(instrument) AS instrument FROM log WHERE (sim_instrument > 0 OR instrument > 0) AND dual > 0 AND (extract(year from age(NOW()::date, date)) * 12 + extract(month from age(NOW()::date, date))) <= 1").Scan(&reqs.RecentInstrumentInstruction)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(reqs)
}
