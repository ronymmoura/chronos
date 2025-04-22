package models

import "time"

type Process struct {
	ID               int32      `db:"id"                 json:"id"`
	Name             string     `db:"name"               json:"name"`
	Description      string     `db:"description"        json:"description"`
	Path             string     `db:"path"               json:"path"`
	ExecuteEverySecs int        `db:"execute_every_secs" json:"execute_every_secs"`
	CreatedAt        time.Time  `db:"created_at"         json:"created_at"`
	LastRun          *time.Time `db:"last_run"           json:"last_run"`
}
