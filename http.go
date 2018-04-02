package main

import (
	"database/sql"
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
	Id     int64  `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
}

func sendNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache,must-revalidate")
}

func respondWithJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	sendNoCacheHeaders(w)

	jsontext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Can't marshal data: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsontext))

	return nil
}

func respondWithJsonAndStatus(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) error {
	sendNoCacheHeaders(w)

	jsontext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Can't marshal data: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
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

func initUser(db *sql.DB, r *http.Request) (*models.User, error) {
	login, err := r.Cookie("Login")
	password, errP := r.Cookie("Password")
	if err != nil || errP != nil {
		return nil, fmt.Errorf("No cookie set")
	}

	return models.LoadUser(db, login.Value, password.Value)
}

func (h *HttpHandlers) GetMatches(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	matches, err := models.LoadMatches(h.Env.DB)
	if err != nil {
		log.Print("Can't load matches: ", err)
		return
	}

	if err := respondWithJson(w, r, matches); err != nil {
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

	m := &models.Match{
		Teams: jsonMatch.Teams,
		Date:  jsonMatch.Date,
	}

	err := models.AddMatch(h.Env.DB, m)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Can't save match: %v", err)
		return
	}

	respondWithJsonAndStatus(w, r, &requestResult{Status: "OK", Id: m.Id}, http.StatusCreated)
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

	respondWithJsonAndStatus(w, r, &requestResult{Status: "OK"}, http.StatusCreated)
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

	respondWithJson(w, r, &requestResult{Status: "OK"})
	log.Printf("Match saved: %+v", jsonMatch)
}

func (h *HttpHandlers) PostLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var jsonUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := processBody(w, r, &jsonUser); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Can't process auth: %v", err)
		return
	}

	if len(jsonUser.Login) == 0 || len(jsonUser.Password) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Login or password is empty")
		return
	}

	_, err := models.CheckCredentials(h.Env.DB, jsonUser.Login, jsonUser.Password)
	if err != nil {
		respondWithJson(w, r, &requestResult{Status: "Failed", Text: "Incorrect user or password"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Login",
		Value:   jsonUser.Login,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "Password",
		Value:   models.GetMD5Hash(jsonUser.Password),
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})

	respondWithJson(w, r, &requestResult{Status: "OK"})
	log.Printf("User authorized: %s", jsonUser.Login)
}

func (h *HttpHandlers) GetLogout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:    "Login",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour * 24),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "Password",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour * 24),
	})
	respondWithJson(w, r, &requestResult{Status: "OK"})
}

func (h *HttpHandlers) GetLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u, err := initUser(h.Env.DB, r)
	if err != nil {
		respondWithJson(w, r, &requestResult{Status: "Fail", Text: err.Error()})
	} else {
		respondWithJson(w, r, &u)
	}
}

func (h *HttpHandlers) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := models.LoadUsers(h.Env.DB)
	if err != nil {
		log.Print("Can't load users: ", err)
		return
	}

	if err := respondWithJson(w, r, users); err != nil {
		log.Print("Can't send response: ", err)
		return
	}
}

func (h *HttpHandlers) GetToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
