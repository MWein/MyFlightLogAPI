package endpoints

import (
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
)

func PlaneImg(w http.ResponseWriter, r *http.Request) {
	idents, ok := r.URL.Query()["ident"]

	if !ok || len(idents[0]) < 1 {
		fmt.Fprintf(w, "ident is required")
		return
	}

	var image []byte
	err := database.DBConnection.QueryRow("SELECT pic FROM plane WHERE ident = $1", idents[0]).Scan(&image)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}
