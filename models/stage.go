package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Stage struct {
	Id        int64     `json:id`
	Name      string    `json:name`
	StartDate time.Time `json:startDate`
	EndDate   time.Time `json:endDate`
}

const (
	CREATE_STAGES_TABLE  = "CREATE TABLE IF NOT EXISTS Stages(name, start_date, end_date)"
	SELECT_ALL_STAGES    = "SELECT rowid, name, start_date, end_date FROM Stages"
	SELECT_CURRENT_STAGE = SELECT_ALL_STAGES + " WHERE date('now') >= start_date AND date('now') <= end_date"
)

func InitStagesTable(db *sql.DB) error {
	_, err := db.Exec(CREATE_STAGES_TABLE)
	return err
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanStage(row scannable) (*Stage, error) {
	var start, end string
	var err error
	s := new(Stage)

	if err := row.Scan(&s.Id, &s.Name, &start, &end); err != nil {
		return nil, fmt.Errorf("Can't extract stage information from database: %s", err.Error())
	}
	if s.StartDate, err = time.Parse(TIMEFORMAT, start); err != nil {
		return nil, fmt.Errorf("Can't parse date %s: %s", start, err.Error())
	}
	if s.EndDate, err = time.Parse(TIMEFORMAT, end); err != nil {
		return nil, fmt.Errorf("Can't parse date %s: %s", end, err.Error())
	}
	return s, nil
}

func GetCurrentStage(db *sql.DB) (*Stage, error) {
	row := db.QueryRow(SELECT_CURRENT_STAGE)

	return scanStage(row)
}

func LoadStages(db *sql.DB) ([]*Stage, error) {
	rows, err := db.Query(SELECT_ALL_STAGES)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stages := make([]*Stage, 0)
	for rows.Next() {
		s, err := scanStage(rows)
		if err != nil {
			return nil, err
		}
		stages = append(stages, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stages, nil
}
