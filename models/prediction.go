package models

import (
	"database/sql"
	"fmt"
)

type Prediction struct {
	UserId  int64  `json:"userId"`
	MatchId int64  `json:"-"`
	Score   string `json:"score"`
}

const (
	ADD                                 = "INSERT INTO Predictions(user_id, match_id, score) VALUES($1,$2,$3)"
	UPDATE                              = "UPDATE Predictions SET score=$1 WHERE user_id=$2 AND match_id=$3"
	SELECT_ALL_PREDICTIONS              = "SELECT user_id, match_id, score FROM Predictions"
	SELECT_PREDICTIONS_IN_MATCHES_RANGE = SELECT_ALL_PREDICTIONS + " WHERE match_id >= ? AND match_id <= ?"
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
	_, err := db.Exec(UPDATE, pred.Score, pred.UserId, pred.MatchId)

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
	return err
}

func loadPredictions(rows *sql.Rows) ([]*Prediction, error) {
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

func LoadPredictionsByMatches(db *sql.DB, matches []*Match) ([]*Prediction, error) {
	if len(matches) == 0 {
		return make([]*Prediction, 0), nil
	}

	maxId := matches[0].Id
	minId := maxId
	for _, m := range matches {
		if minId > m.Id {
			minId = m.Id
			continue
		}

		if maxId < m.Id {
			maxId = m.Id
		}
	}

	rows, err := db.Query(SELECT_PREDICTIONS_IN_MATCHES_RANGE, minId, maxId)
	if err != nil {
		return nil, fmt.Errorf("Can't load predictions: %s", err.Error())
	}

	defer rows.Close()

	return loadPredictions(rows)
}
