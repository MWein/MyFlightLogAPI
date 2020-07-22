package endpoints

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func PlaneImg(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

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

	elapsed := time.Since(start)
	fmt.Printf("Connection took %s\n", elapsed)

	start = time.Now()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	fmt.Printf("Ping took %s\n", elapsed)

	start = time.Now()
	var image []byte
	err = db.QueryRow("SELECT pic FROM plane WHERE ident = $1", idents[0]).Scan(&image)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}
	elapsed = time.Since(start)
	fmt.Printf("Getting image took %s\n", elapsed)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}
