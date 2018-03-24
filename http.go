package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aelnor/vangothrone/config"
	"github.com/aelnor/vangothrone/models"
)

type HttpHandlers struct {
	Env *config.Env
}

type requestResult struct {
	Status string `json:"status"`
	Text   string `json:"text,omitempty"`
}

func sendNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func respondWithJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	sendNoCacheHeaders(w)

	jsontext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Can't marshal data: %s", err.Error())
	}

	fmt.Fprintf(w, string(jsontext))

	return nil
}

func processBody(w http.ResponseWriter, r *http.Request, result interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Can't read body from request: %v", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("Can't parse requst body: %v", err)
	}
	return nil
}

func getMatches(env *config.Env, w http.ResponseWriter, r *http.Request) {
	matches, err := models.LoadMatches(env.DB)
	if err != nil {
		log.Print("Can't load matches: ", err)
		return
	}

	if err := respondWithJson(w, matches); err != nil {
		log.Print("Can't send response: ", err)
		return
	}
}

func postMatches(env *config.Env, w http.ResponseWriter, r *http.Request) {
	var jsonMatch struct {
		Teams [2]string `json:"teams"`
		Date  time.Time `json:"date"`
	}

	if err := processBody(w, r, &jsonMatch); err != nil {
		log.Printf("Can't process match adding: %v", err)
		return
	}

	err := models.AddMatch(env.DB, &models.Match{
		Teams: jsonMatch.Teams,
		Date:  jsonMatch.Date,
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Can't save match: %v", err)
		return
	}

	respondWithJson(w, &requestResult{Status: "OK"})
	log.Printf("Match added: %+v", jsonMatch)
}

func (h *HttpHandlers) Matches(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		getMatches(h.Env, w, r)
	case r.Method == http.MethodPost:
		postMatches(h.Env, w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Printf("Query to /matches with %s from %s", r.Method, r.RemoteAddr)
		return
	}
}

func (h *HttpHandlers) PutPrediction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		log.Printf("Query to /prediction with %s from %s", r.Method, r.RemoteAddr)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Can't read prediction from request")
		return
	}

	var jsonPrediction struct {
		UserId  int64  `json:"userId"`
		MatchId int64  `json:"matchId"`
		Score   string `json:"score"`
	}

	err = json.Unmarshal(body, &jsonPrediction)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Can't parse prediction: %v", err)
		return
	}

	err = models.SavePrediction(h.Env.DB, &models.Prediction{
		UserId:  jsonPrediction.UserId,
		MatchId: jsonPrediction.MatchId,
		Score:   jsonPrediction.Score,
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Can't save prediction: %v", err)
		return
	}

	respondWithJson(w, &requestResult{Status: "OK"})
	log.Printf("Saved prediction: %+v", jsonPrediction)
}
