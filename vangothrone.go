package main

import (
	"encoding/json"
	"fmt"
	"github.com/aelnor/vangothrone/config"
	"github.com/aelnor/vangothrone/models"
	"log"
	"net/http"
)

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

func InitEnvironment() (*config.Env, error) {
	db, err := config.InitDatabase()
	if err != nil {
		return nil, fmt.Errorf("Can't init database: %s", err.Error())
	}
	if err := models.InitUsersTable(db); err != nil {
		return nil, fmt.Errorf("Can't init database 'Users': %s", err.Error())
	}
	if err := models.InitMatchesTable(db); err != nil {
		return nil, fmt.Errorf("Can't init database 'Matches': %s", err.Error())
	}
	if err := models.InitPredictionsTable(db); err != nil {
		return nil, fmt.Errorf("Can't init database 'Predictions': %s", err.Error())
	}
	log.Printf("Database Initialized")

	env := &config.Env{
		DB: db,
	}
	return env, nil
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	jsontext, err := json.Marshal(models.Teams)

	if err != nil {
		fmt.Printf("Can't marshal teams: %v", err)
		return
	}

	fmt.Fprintf(w, string(jsontext))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	}
}

func main() {
	env, err := InitEnvironment()
	if err != nil {
		log.Fatal("Can't init environment: ", err)
	}
	hh := &HttpHandlers{Env: env}
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/teams", teamsHandler)
	http.HandleFunc("/matches", hh.Matches)
	http.HandleFunc("/prediction", hh.PutPrediction)
	log.Printf("Preparations finished, serving")
	log.Fatal(http.ListenAndServe(":8383", nil))
}
