package repositories

import (
	"time"

	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/models"
)

type ProcessRepository struct {
	DB *db.DB
}

func NewProcessRepository(db *db.DB) *ProcessRepository {
	return &ProcessRepository{
		DB: db,
	}
}

//-------------------------
// Insert
//-------------------------

const insertQuery = `
INSERT INTO process
(
	name,
	description,
	path,
	execute_every_secs
) VALUES (
	:name,
	:description,
	:path,
	:execute_every_secs
) RETURNING *
`

func (r *ProcessRepository) Insert(args *models.Process) (*models.Process, error) {
	rows, err := r.DB.Conn.NamedQuery(insertQuery, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	process := &models.Process{}
	for rows.Next() {
		err := rows.StructScan(&process)
		if err != nil {
			return nil, err
		}
	}

	return process, nil
}

//--------------------------------
// Select All
//--------------------------------

const selectAllQuery = `SELECT * FROM process`

func (r *ProcessRepository) SelectAll() ([]*models.Process, error) {
	rows := []*models.Process{}
	if err := r.DB.Conn.Select(&rows, selectAllQuery); err != nil {
		return nil, err
	}

	return rows, nil
}

//----------------------------------
// Update status
//----------------------------------

const updateStatusQuery = `
	UPDATE process
	SET last_run = $2
	WHERE id = $1
	`

func (r *ProcessRepository) UpdateStatus(id int32, lastRun time.Time) error {
	if _, err := r.DB.Conn.Exec(updateStatusQuery, id, lastRun); err != nil {
		return err
	}

	return nil
}
