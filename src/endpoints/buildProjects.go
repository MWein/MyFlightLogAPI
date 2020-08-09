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
	rows, _ := database.DBConnection.Query("SELECT id, name FROM build")

	var projects Projects
	for rows.Next() {
		var project Project
		rows.Scan(&project.Id, &project.Name)

		// Placeholder data
		project.Hours = 140.15
		project.LastEntry = "2020-01-03"

		projects = append(projects, project)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(projects)
}
