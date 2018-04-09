package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Match struct {
	Id          int64         `json:"id"`
	Teams       [2]string     `json:"teams"`
	Date        time.Time     `json:"date"`
	Result      string        `json:"result"`
	Predictions []*Prediction `json:"predictions"`
}

const (
	TIMEFORMAT           = "2006-01-02T15:04:05Z0700"
	CREATE_MATCHES_TABLE = "CREATE TABLE IF NOT EXISTS Matches(team_a, team_b, date, result)"
	SELECT_ALL_MATCHES   = "SELECT rowid, team_a, team_b, date, result FROM Matches ORDER BY date ASC"
	SELECT_MATCH_BY_ID   = "SELECT rowid, team_a, team_b, date, result FROM Matches WHERE rowid=?"
)

func InitMatchesTable(db *sql.DB) error {
	_, err := db.Exec(CREATE_MATCHES_TABLE)
	return err
}

func AddMatch(db *sql.DB, m *Match) error {
	if len(m.Teams[0]) == 0 || len(m.Teams[1]) == 0 {
		return fmt.Errorf("There should be 2 teams")
	}
	date := m.Date.UTC().Format(TIMEFORMAT)
	result, err := db.Exec("INSERT INTO Matches(team_a, team_b, date, result) VALUES(?,?,?,?)", m.Teams[0], m.Teams[1], date, m.Result)

	if err == nil {
		m.Id, _ = result.LastInsertId()
	}
	return err
}

func SaveMatch(db *sql.DB, m *Match) error {
	if m.Id == 0 {
		return fmt.Errorf("MatchId is null")
	}

	fields := make([]string, 0, 4)
	values := make([]string, 0, 4)

	if len(m.Teams[0]) != 0 && len(m.Teams[1]) != 0 {
		fields = append(fields, "team_a")
		values = append(values, m.Teams[0])
		fields = append(fields, "team_b")
		values = append(values, m.Teams[1])
	}

	if !m.Date.IsZero() {
		fields = append(fields, "date")
		values = append(values, m.Date.Format(TIMEFORMAT))
	}

	if len(m.Result) != 0 {
		fields = append(fields, "result")
		values = append(values, m.Result)
	}

	query := "UPDATE Matches SET "
	for i, val := range fields {
		query += val + "='" + values[i] + "'"
		if i != len(fields)-1 {
			query += ", "
		}
	}

	query += " WHERE rowid=?"

	res, err := db.Exec(query, m.Id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err

	}
	if rows != 1 {
		return fmt.Errorf("No such match")
	}
	return nil
}

func LoadMatches(db *sql.DB) ([]*Match, error) {
	rows, err := db.Query(SELECT_ALL_MATCHES)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := make([]*Match, 0)

	for rows.Next() {
		var id int64
		var teamA, teamB, date, result string
		if err := rows.Scan(&id, &teamA, &teamB, &date, &result); err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(TIMEFORMAT, date)
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

func LoadMatch(db *sql.DB, id int64) (*Match, error) {
	row := db.QueryRow(SELECT_MATCH_BY_ID, id)

	m := new(Match)
	var date string

	err := row.Scan(&m.Id, &m.Teams[0], &m.Teams[1], &date, &m.Result)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("No match with id: ", id)
	case err != nil:
		return nil, err
	}

	m.Date, err = time.Parse(TIMEFORMAT, date)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Match) IsStarted() bool {
	return m.Date.After(time.Now().UTC())
}
