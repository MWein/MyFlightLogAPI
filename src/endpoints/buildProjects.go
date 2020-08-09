package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

type Project struct {
	Id string  `json:"id"`
	Name  string  `json:"name"`
}

type Projects []Project

func BuildProjects(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DBConnection.Query("SELECT id, name FROM build")

	var projects Projects
	for rows.Next() {
		var project Project
		rows.Scan(&project.Id, &project.Name)

		projects = append(projects, project)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(projects)
}
