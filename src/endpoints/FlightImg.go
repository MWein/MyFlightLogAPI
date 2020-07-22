package endpoints

import (
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

func FlightImg(w http.ResponseWriter, r *http.Request) {
	fmt.Println(database.DBConnection)


	id, ok := r.URL.Query()["imgid"]

	if !ok || len(id[0]) < 1 {
		fmt.Fprintf(w, "Image ID is required")
		return
	}

	var image []byte
	err := database.DBConnection.QueryRow("SELECT data FROM pictures WHERE id = $1", id[0]).Scan(&image)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}
