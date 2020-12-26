package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var DBConnection *sql.DB
var Cache *cache.Cache

func StartDB() {
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

	Cache = cache.New(5*time.Minute, 10*time.Minute)

	fmt.Println("Cache Ready")
}
