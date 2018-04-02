package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Prediction struct {
	UserId  int64  `json:"userId"`
	MatchId int64  `json:"-"`
	Score   string `json:"score"`
}

const (
	ADD             = "INSERT INTO Predictions(user_id, match_id, score) VALUES($1,$2,$3)"
	UPDATE          = "UPDATE Predictions SET score=$3 WHERE user_id=$1 AND match_id=$2"
	FIND_BY_MATCHES = "SELECT user_id, match_id, score FROM Predictions WHERE match_id IN "
)

func InitPredictionsTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS Predictions(user_id, match_id, score)")

	return err
}

func createPrediction(db *sql.DB, pred *Prediction) error {
	_, err := db.Exec(ADD, pred.UserId, pred.MatchId, pred.Score)

	return err
}

func updatePrediction(db *sql.DB, pred *Prediction) error {
	_, err := db.Exec(UPDATE, pred.UserId, pred.MatchId, pred.Score)

	return err
}

func SavePrediction(db *sql.DB, pred *Prediction) error {
	if pred.MatchId == 0 {
		return fmt.Errorf("Match ID is null")
	}
	if pred.UserId == 0 {
		return fmt.Errorf("User ID is null")
	}

	row := db.QueryRow("SELECT rowid, score FROM Predictions WHERE user_id=? AND match_id=?", pred.UserId, pred.MatchId)

	var score string
	var id int
	err := row.Scan(&id, &score)

	switch {
	case err == nil:
		if score != pred.Score {
			err = updatePrediction(db, pred)
		}
	case err == sql.ErrNoRows:
		err = createPrediction(db, pred)

	}
	if err == nil {
		invalidateMatchesCache()
	}
	return err
}

func LoadPredictions(db *sql.DB, matches []int64) ([]*Prediction, error) {
	query := FIND_BY_MATCHES + "(" + strings.Join(int64ToString(matches), ",") + ")"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Can't load predictions by matches: %s", err.Error())
	}

	defer rows.Close()

	result := make([]*Prediction, 0)
	for rows.Next() {
		pred := new(Prediction)
		rows.Scan(&pred.UserId, &pred.MatchId, &pred.Score)
		result = append(result, pred)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Can't parse predictions: %s", rows.Err())
	}

	return result, nil
}

func int64ToString(input []int64) []string {
	output := make([]string, len(input))
	for i, el := range input {
		output[i] = strconv.FormatInt(el, 10)
	}
	return output
}
