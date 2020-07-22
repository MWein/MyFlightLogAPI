package database

import (
	"database/sql"
	"fmt"
)

var DBConnection *sql.DB

func StartDB() {
	fmt.Println(DBConnection)

	// Spin up the database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DBConnection, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = DBConnection.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database Ready")

	fmt.Println(DBConnection)
}
