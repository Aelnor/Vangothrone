package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

func matchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	matches, err := LoadMatches(db)
	if err != nil {
		log.Print("Can't load matches: ", err)
	}
	jsontext, err := json.MarshalIndent(matches, "", "  ")

	if err != nil {
		fmt.Printf("Can't marshal matches: %s", err.Error())
	}

	fmt.Fprintf(w, string(jsontext))
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	jsontext, err := json.Marshal(Teams)

	if err != nil {
		fmt.Printf("Can't marshal teams: %s", err.Error())
	}

	fmt.Fprintf(w, string(jsontext))
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
	if err := initDatabase(); err != nil {
		log.Fatal("Can't init database: ", err)
	}
	defer db.Close()

	if err := initUsersTable(db); err != nil {
		log.Fatal("Can't init database 'Users': ", err)
	}
	if err := initMatchesTable(db); err != nil {
		log.Fatal("Can't init database 'Matches': ", err)
	}
	log.Printf("Database Initialized")
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/teams", teamsHandler)
	http.HandleFunc("/", matchesHandler)
	log.Printf("Preparations finished, serving")
	log.Fatal(http.ListenAndServe(":8383", nil))
}
