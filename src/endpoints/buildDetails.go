package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type BuildEntry struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Minutes     int      `json:"minutes"`
	Rivets      int      `json:"rivets"`
	Description string   `json:"description"`
	Pictures    []string `json:"pictures"`
	Phase       string   `json:"phase"`
}
type BuildEntries []BuildEntry

type Expense struct {
	Date        string  `json:"date"`
	Cost        float32 `json:"cost"`
	Projected   bool    `json:"projected"`
	Description string  `json:"description"`
}
type Expenses []Expense

type Phase struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Complete bool         `json:"complete"`
	Expenses Expenses     `json:"expenses"`
	Entries  BuildEntries `json:"entries"`
}
type Phases []Phase

type Build struct {
	Name   string `json:"name"`
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
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	rows, _ := database.DBConnection.Query(`SELECT id, name, complete FROM build_phase WHERE build_id = $1`, buildId[0])
	phases := Phases{}
	for rows.Next() {
		var phase Phase
		rows.Scan(&phase.Id, &phase.Name, &phase.Complete)

		// Get entries
		rows, _ := database.DBConnection.Query(`SELECT title, date, minutes, rivets, description, (SELECT name FROM build_phase WHERE id = $1) AS phase, id FROM build_log WHERE phase_id = $1 ORDER BY date DESC`, phase.Id)
		phase.Entries = BuildEntries{}
		for rows.Next() {
			var entry BuildEntry
			var buildLogId string
			rows.Scan(&entry.Title, &entry.Date, &entry.Minutes, &entry.Rivets, &entry.Description, &entry.Phase, &buildLogId)

			// Remove timestamp from date
			entry.Date = entry.Date[0:10]

			// Get pictures
			pictureIdRows, _ := database.DBConnection.Query(`SELECT id FROM build_log_picture WHERE buildlogid = $1`, buildLogId)
			entry.Pictures = []string{}
			for pictureIdRows.Next() {
				var picId string
				pictureIdRows.Scan(&picId)
				entry.Pictures = append(entry.Pictures, picId)
			}

			phase.Entries = append(phase.Entries, entry)
		}

		// Get expenses
		rows, _ = database.DBConnection.Query(`SELECT description, date, cost, projected FROM build_expense WHERE phase_id = $1`, phase.Id)
		phase.Expenses = Expenses{}
		for rows.Next() {
			var expense Expense
			rows.Scan(&expense.Description, &expense.Date, &expense.Cost, &expense.Projected)

			// Remove timestamp from date
			if (expense.Date != "") {
				expense.Date = expense.Date[0:10]
			}

			phase.Expenses = append(phase.Expenses, expense)
		}

		phases = append(phases, phase)
	}

	build.Phases = phases

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(build)
}
