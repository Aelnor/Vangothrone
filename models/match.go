package models

import (
	"database/sql"
	"time"
)

type Match struct {
	Id          int          `json:"id"`
	Teams       [2]string    `json:"teams"`
	Date        time.Time    `json:"date"`
	Result      string       `json:"result"`
	Predictions []Prediction `json:"predictions"`
}

const timeformat = "Jan _2 2006 15:04"

func InitMatchesTable(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Matches(team_a, team_b, date, result)"); err != nil {
		return err
	}

	return addMatch(db, Match{Id: 0, Teams: [2]string{"DAL", "SEO"}, Date: time.Now(), Result: "0:4"})
}

func addMatch(db *sql.DB, m Match) error {
	statement, err := db.Prepare("INSERT INTO Matches(team_a, team_b, date, result) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	date := m.Date.Format(timeformat)
	_, err = statement.Exec(m.Teams[0], m.Teams[1], date, m.Result)

	return err
}

func LoadMatches(db *sql.DB) ([]*Match, error) {
	rows, err := db.Query("SELECT rowid, team_a, team_b, date, result FROM Matches ORDER BY date ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := make([]*Match, 0)

	for rows.Next() {
		var id int
		var teamA, teamB, date, result string
		if err := rows.Scan(&id, &teamA, &teamB, &date, &result); err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(timeformat, date)
		if err != nil {
			return nil, err
		}
		matches = append(matches, &Match{Id: id, Teams: [2]string{teamA, teamB}, Date: parsedDate, Result: result})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return matches, nil
}
