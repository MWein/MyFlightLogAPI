package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type Project struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Hours     float32 `json:"hours"`
	LastEntry string  `json:"lastEntry"`
}

type Projects []Project

func BuildProjects(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DBConnection.Query(`
	SELECT
		id, name,
		(SELECT ROUND(sum(minutes) / 60, 2) FROM build_log WHERE phase_id = ANY (SELECT id FROM build_phase WHERE build_id = build.id)) AS hours,
		(SELECT date FROM build_log WHERE phase_id = ANY (SELECT id FROM build_phase WHERE build_id = build.id) ORDER BY date DESC LIMIT 1) AS last_entry
	FROM build`)

	var projects Projects
	for rows.Next() {
		var project Project
		rows.Scan(&project.Id, &project.Name, &project.Hours, &project.LastEntry)

		projects = append(projects, project)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(projects)
}
