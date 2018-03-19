package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, `[
  { "id": 1,
    "teams": ["VAL", "NYE"],
    "date": "12 Mar 2018, 02:00",
    "result": "3:1",
    "resultFrozen": "true",
    "predictions": [{"userId": 1, "score": "3:1", "matchId": 1},
                  {"userId": 2, "score": "3:2", "matchId": 1}]},
  { "id": "2",
    "teams": ["VAL", "NYE"],
    "date": "13 Mar 2018, 02:00",
    "result": "1:3",
    "resultFrozen": false,
    "predictions": [{"userId": 1, "score": "0:4", "matchId": 2},
                  {"userId": 2, "score": "2:3", "matchId": 2}]},
  { "id": 3,
    "teams": ["SEO", "NYE"],
    "date": "19 Mar 2018, 02:00",
    "result": "",
    "resultFrozen": false,
    "predictions": [{"userId": 1, "score": "3:1", "matchId": 3},
                  {"userId": 2, "score": "3:2", "matchId": 3}]},
  { "id": 4,
    "teams": ["VAL", "SEO"],
    "date": "20 Mar 2018, 02:00",
    "result": "",
    "resultFrozen": false,
    "predictions": []}
]`)
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json, err := json.Marshal(Teams)

	if err != nil {
		fmt.Printf("Can't marshal teams: %s", err.Error())
	}

	fmt.Fprintf(w, string(json))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	}
	json, err := json.Marshal(CurrentUser)

	if err != nil {
		fmt.Printf("Can't marshal current user: %s", err.Error())
	}

	fmt.Fprintf(w, string(json))

}

func main() {
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/teams", teamsHandler)
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(":8383", nil))
}
