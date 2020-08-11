package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)


type BuildEntry struct {
	Date       string `json:"date"`
	Minutes     int `json:"minutes"`
	Rivets int   `json:"rivets"`
	Description string  `json:"description"`
	Pictures []string `json:"pictures"`
}
type BuildEntries []BuildEntry


type Expense struct {
	Date       string `json:"date"`
	Cost     float32 `json:"cost"`
	Projected bool   `json:"projected"`
	Description string  `json:"description"`
}
type Expenses []Expense


type Phase struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Expenses Expenses  `json:"expenses"`
	Entries BuildEntries  `json:"entries"`
}
type Phases []Phase


type Build struct {
	Name string `json:"name"`
	Phases Phases `json:"phases"`
}




func BuildDetails(w http.ResponseWriter, r *http.Request) {
	buildId, ok := r.URL.Query()["buildid"]

	if !ok || len(buildId[0]) < 1 {
		fmt.Fprintf(w, "Build ID is required")
		return
	}


	var build Build
	err := database.DBConnection.QueryRow("SELECT name FROM build WHERE id = $1", buildId[0]).Scan(&build.Name)
	if (err != nil) {
		fmt.Fprintf(w, "Not Found")
		return
	}


	rows, _ := database.DBConnection.Query(`SELECT id, name, complete FROM build_phase WHERE build_id = $1`, buildId[0])
	var phases Phases
	for rows.Next() {
		var phase Phase
		rows.Scan(&phase.Id, &phase.Name, &phase.Complete)

		// Get entries
		rows, _ := database.DBConnection.Query(`SELECT date, minutes, rivets, description FROM build_log WHERE phase_id = $1`, phase.Id)
		phase.Entries = BuildEntries{}
		for rows.Next() {
			var entry BuildEntry
			rows.Scan(&entry.Date, &entry.Minutes, &entry.Rivets, &entry.Description)

			// Remove timestamp from date
			entry.Date = entry.Date[0:10]

			// TODO Get pictures
			entry.Pictures = []string{}

			phase.Entries = append(phase.Entries, entry)
		}

		// TODO Get expenses
		phase.Expenses = Expenses{}

		phases = append(phases, phase)
	}


	build.Phases = phases


	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(build)
}
