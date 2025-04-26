package models

import "time"

type Process struct {
	ID               int32     `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	Path             string    `db:"path" json:"-"`
	Env              string    `db:"env" json:"env"`
	Status           string    `db:"status" json:"status"`
	ExecuteEverySecs int       `db:"execute_every_secs" json:"execute_every_secs"`
	Running          bool      `db:"running" json:"running"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}
