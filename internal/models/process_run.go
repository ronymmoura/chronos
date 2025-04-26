package models

import "time"

type ProcessRun struct {
	ID        int32     `db:"id" json:"id"`
	ProcessID int32     `db:"process_id" json:"process_id"`
	StartedAt time.Time `db:"started_at" json:"started_at"`
	EndedAt   time.Time `db:"ended_at" json:"ended_at"`
	Success   bool      `db:"success" json:"success"`

	Process *Process `db:"process" json:"process"`

	Logs []ProcessRunLog `json:"process_run_logs"`
}
