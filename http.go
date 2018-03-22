package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aelnor/vangothrone/config"
)

type HttpHandlers struct {
	Env *config.Env
}

func respondWithJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsontext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Can't marshal matches: %s", err.Error())
	}

	fmt.Fprintf(w, string(jsontext))
}

func (h *HttpHandlers) AllMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := models.LoadMatches(h.Env.DB)
	if err != nil {
		log.Print("Can't load matches: ", err)
	}

	if err := respondWithJson(w, matches); err != nil {
		log.Print("Can't create response: ", err)
	}
}
