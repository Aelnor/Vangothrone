package main

import "time"

type Match struct {
	id          int
	teams       []string
	date        time.Time
	result      string
	predictions []Prediction
}
