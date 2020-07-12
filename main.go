package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/MWein/MyFlightLogAPI/src/endpoints"
)


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func otherPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Other Page Endpoint Hit")
}


func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/other", otherPage)

	http.HandleFunc("/ping", endpoints.Ping)
	http.HandleFunc("/log", endpoints.FlightLogs)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}