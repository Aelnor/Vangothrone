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
	SELECT_CURRENT_STAGE = "SELECT rowid, name, start_date, end_date FROM Stages WHERE date('now') >= start_date AND date('now') <= end_date"
)

func InitStagesTable(db *sql.DB) error {
	_, err := db.Exec(CREATE_STAGES_TABLE)
	return err
}

func GetCurrentStage(db *sql.DB) (*Stage, error) {
	row := db.QueryRow(SELECT_CURRENT_STAGE)
	s := new(Stage)
	var start, end string
	err := row.Scan(&s.Id, &s.Name, &start, &end)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("No current stage found")
	case err != nil:
		return nil, err
	}

	if s.StartDate, err = time.Parse(TIMEFORMAT, start); err != nil {
		return nil, fmt.Errorf("Can't parse date %s: %s", start, err.Error())
	}
	if s.EndDate, err = time.Parse(TIMEFORMAT, end); err != nil {
		return nil, fmt.Errorf("Can't parse date %s: %s", end, err.Error())
	}

	return s, nil

}
