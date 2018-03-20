package main

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

const Timeformat = "Jan _2 2006 15:04"

func initMatchesTable(db *sql.DB) error {
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

	date := m.Date.Format(Timeformat)
	_, err = statement.Exec(m.Teams[0], m.Teams[1], date, m.Result)

	return err
}

func extend(slice []Match, element Match) []Match {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]Match, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func LoadMatches(db *sql.DB) ([]Match, error) {
	rows, err := db.Query("SELECT rowid, team_a, team_b, date, result FROM Matches ORDER BY date ASC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var matches []Match

	for rows.Next() {
		var id int
		var teamA, teamB, date, result string
		if err := rows.Scan(&id, &teamA, &teamB, &date, &result); err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(Timeformat, date)
		if err != nil {
			return nil, err
		}
		matches = extend(matches, Match{Id: id, Teams: [2]string{teamA, teamB}, Date: parsedDate, Result: result})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return matches, nil
}
