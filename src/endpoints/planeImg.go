package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"
)

func PlaneImg(w http.ResponseWriter, r *http.Request) {
	idents, ok := r.URL.Query()["ident"]

	if !ok || len(idents[0]) < 1 {
		fmt.Fprintf(w, "ident is required")
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

	planeImgStatement := `SELECT pic FROM plane WHERE ident = $1`

	var image []byte
	err = db.QueryRow(planeImgStatement, idents[0]).Scan(&image)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}
