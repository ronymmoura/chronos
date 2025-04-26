package repositories

import (
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
	env,
	execute_every_secs
) VALUES (
	:name,
	:description,
	:path,
	:env,
	:execute_every_secs
)
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

// --------------------------------
// Select By ID
// --------------------------------

const selectByIDQuery = `SELECT * FROM process WHERE id = $1`

func (r *ProcessRepository) GetByID(id int32) (*models.Process, error) {
	process := &models.Process{}
	if err := r.DB.Conn.Get(process, selectByIDQuery, id); err != nil {
		return nil, err
	}

	return process, nil
}

// --------------------------------
// Select Actives
// --------------------------------

const selectActivesQuery = `SELECT * FROM process WHERE status = 'active'`

func (r *ProcessRepository) SelectActives() ([]*models.Process, error) {
	rows := []*models.Process{}
	if err := r.DB.Conn.Select(&rows, selectActivesQuery); err != nil {
		return nil, err
	}

	return rows, nil
}

//----------------------------------
// Update runnint
//----------------------------------

const updateRunningQuery = `
	UPDATE process
	SET running = $2
	WHERE id = $1
	`

func (r *ProcessRepository) UpdateRunning(id int32, running bool) error {
	if _, err := r.DB.Conn.Exec(updateRunningQuery, id, running); err != nil {
		return err
	}

	return nil
}
