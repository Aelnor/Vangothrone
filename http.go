package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aelnor/vangothrone/config"
	"github.com/aelnor/vangothrone/models"
	"github.com/julienschmidt/httprouter"
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
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

func (h *HttpHandlers) GetMatches(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	matches, err := models.LoadMatches(h.Env.DB)
	if err != nil {
		log.Print("Can't load matches: ", err)
		return
	}

	if err := respondWithJson(w, matches); err != nil {
		log.Print("Can't send response: ", err)
		return
	}
}

func (h *HttpHandlers) PostMatches(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var jsonMatch struct {
		Teams [2]string `json:"teams"`
		Date  time.Time `json:"date"`
	}

	if err := processBody(w, r, &jsonMatch); err != nil {
		log.Printf("Can't process match adding: %v", err)
		return
	}

	err := models.AddMatch(h.Env.DB, &models.Match{
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

func (h *HttpHandlers) PutPredictions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func (h *HttpHandlers) PutMatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	paramId := p.ByName("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Bad match id: %s", paramId)
		return
	}

	var jsonMatch struct {
		Teams  [2]string `json:"teams"`
		Date   time.Time `json:"date"`
		Result string    `json:"result"`
	}

	if err = processBody(w, r, &jsonMatch); err != nil {
		log.Printf("Can't process match editing: %v", err)
		return
	}

	err = models.SaveMatch(h.Env.DB, &models.Match{Id: id, Teams: jsonMatch.Teams, Date: jsonMatch.Date, Result: jsonMatch.Result})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Can't save match: %v", err)
		return
	}

	respondWithJson(w, &requestResult{Status: "OK"})
	log.Printf("Match saved: %+v", jsonMatch)
}
