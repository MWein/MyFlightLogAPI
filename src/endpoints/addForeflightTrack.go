package endpoints

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/MWein/MyFlightLogAPI/src/database"
	"github.com/MWein/MyFlightLogAPI/src/utils"
)

func AddForeflightTrack(w http.ResponseWriter, r *http.Request) {
	foreflightId, ok := r.URL.Query()["foreflightId"]

	// TODO Check if ID is in proper format

	if !ok || len(foreflightId[0]) < 1 {
		fmt.Fprintf(w, "foreflightId is required")
		return
	}

	var flightId string
	flightIdQuery, flightIdProvided := r.URL.Query()["flightId"]

	if flightIdProvided {
		// Verify the flight ID exists in the database
		var exists bool
		database.DBConnection.QueryRow("SELECT count(*) = 1 FROM log WHERE id = $1", flightIdQuery[0]).Scan(&exists)

		fmt.Println(exists)

		if !exists {
			fmt.Fprintf(w, "Flight ID is invalid")
			return
		}

		// Set the flight ID to the URL Query
		flightId = flightIdQuery[0]
	}

	// Request CSV from Foreflight
	csvUrl := fmt.Sprintf("https://plan.foreflight.com/tracklogs/export/%s/csv", foreflightId[0])
	resp, _ := http.Get(csvUrl)

	// Handle Errors
	if resp.Status == "404 Not Found" {
		fmt.Fprintf(w, "The foreflight log %s was not found", foreflightId[0])
		return
	} else if resp.Status != "200 OK" {
		fmt.Fprintf(w, "Something went wrong")
		return
	}

	// Read the CSV and load
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','

	var data [][]string
	for {
		nextLine, err := reader.Read()

		if err == io.EOF {
			break
		}

		data = append(data, nextLine)
	}

	// Get the date of the flight
	timestampFl, _ := strconv.ParseFloat(data[3][0], 64)
	timestamp := int64(timestampFl)
	dateStr := time.Unix(timestamp, 0).Format("2006-01-02")

	// If no flight ID was provided
	if !flightIdProvided || len(flightIdQuery[0]) < 1 {
		// Check the database for flights matching the date
		var flightIds []string
		rows, _ := database.DBConnection.Query("SELECT id FROM log WHERE date = $1", dateStr)
		for rows.Next() {
			var id string
			rows.Scan(&id)
			flightIds = append(flightIds, id)
		}

		// If there are no flights, return
		if len(flightIds) == 0 {
			fmt.Fprintf(w, "No logged flights that occured on %s", dateStr)
			return
		}

		// If there are multiple flights, tell the user we need an ID
		if len(flightIds) > 1 {
			multipleFlightError := fmt.Sprintf("Multiple logged flights occured on %s. Please choose a flight ID.\n", dateStr)
			for _, id := range flightIds {
				multipleFlightError = fmt.Sprintf("%s\n%s", multipleFlightError, id)
			}
			fmt.Fprintf(w, multipleFlightError)
			return
		}

		flightId = flightIds[0]
	}

	// Clean out any previous foreflight logs
	database.DBConnection.Exec("DELETE FROM foreflight WHERE flightid = $1", flightId)

	var foreflightTrack [][2]float64

	for index, row := range data {
		// Skip the first 3 lines, go right to the actual flight log data
		if index <= 2 {
			continue
		}

		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println("Invalid latitude format. Skipping.")
			continue
		}

		long, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			fmt.Println("Invalid longitude format. Skipping.")
			continue
		}

		meters := utils.LatLong2Meters(lat, long)

		foreflightTrack = append(foreflightTrack, meters)
	}

	foreflightTrackJSON, _ := json.Marshal(foreflightTrack)
	database.DBConnection.Exec("INSERT INTO foreflight (flightid, data) VALUES ($1, $2)", flightId, foreflightTrackJSON)

	fmt.Fprintf(w, "Saved")
}
