package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"
)

func FlightImg(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["imgid"]

	if !ok || len(id[0]) < 1 {
		fmt.Fprintf(w, "Image ID is required")
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

	var image []byte
	err = db.QueryRow("SELECT data FROM pictures WHERE id = $1", id[0]).Scan(&image)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}
