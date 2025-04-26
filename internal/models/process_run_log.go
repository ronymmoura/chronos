package models

import "time"

type ProcessRunLog struct {
	ID           int32     `db:"id" json:"id"`
	ProcessRunID int32     `db:"process_run_id" json:"process_run_id"`
	Message      string    `db:"message" json:"message"`
	Type         string    `db:"type" json:"type"`
	LogTime      time.Time `db:"log_time" json:"log_time"`
}
