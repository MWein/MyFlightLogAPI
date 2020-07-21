package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func ForeflightTrack(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["flightid"]

	if !ok || len(id[0]) < 1 {
		fmt.Fprintf(w, "Flight ID is required")
		return
	}

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

	rows, err := db.Query("SELECT lat, long FROM foreflight WHERE flightid = $1 ORDER BY sequence", id[0])
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	var foreflightTrack [][2]float64
	for rows.Next() {
		var latLon [2]float64
		rows.Scan(&latLon[0], &latLon[1])

		foreflightTrack = append(foreflightTrack, latLon)
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(foreflightTrack)
}
