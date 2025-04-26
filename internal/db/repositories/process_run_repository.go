package repositories

import (
	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/models"
)

type ProcessRunRepository struct {
	DB *db.DB
}

func NewProcessRunRepository(db *db.DB) *ProcessRunRepository {
	return &ProcessRunRepository{
		DB: db,
	}
}

//--------------------------------
// Insert
//--------------------------------

const insertProcessRunQuery = `
INSERT INTO process_run
(
	process_id,
	started_at,
	ended_at,
	success
) VALUES (
	:process_id,
	:started_at,
	:ended_at,
	:success
);
SELECT SCOPE_IDENTITY()
`

func (r *ProcessRunRepository) Insert(args *models.ProcessRun) (int32, error) {
	rows, err := r.DB.Conn.NamedQuery(insertProcessRunQuery, args)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int32
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// --------------------------------
// Select All
// --------------------------------
const selectAllProcessRunQuery = `SELECT * FROM process_run`

func (r *ProcessRunRepository) SelectAll() ([]*models.ProcessRun, error) {
	rows := []*models.ProcessRun{}
	if err := r.DB.Conn.Select(&rows, selectAllProcessRunQuery); err != nil {
		return nil, err
	}

	return rows, nil
}

//--------------------------------
// Select By Process ID
//--------------------------------

const selectByProcessIDQuery = `SELECT * FROM process_run WHERE process_id = $1 ORDER BY start_at`

func (r *ProcessRunRepository) SelectByProcessID(processID int32) ([]*models.ProcessRun, error) {
	rows := []*models.ProcessRun{}
	if err := r.DB.Conn.Get(&rows, selectByProcessIDQuery, processID); err != nil {
		return nil, err
	}

	return rows, nil
}
