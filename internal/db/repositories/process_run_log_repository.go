package repositories

import (
	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/models"
)

type ProcessRunLogRepository struct {
	DB *db.DB
}

func NewProcessRunLogRepository(db *db.DB) *ProcessRunLogRepository {
	return &ProcessRunLogRepository{
		DB: db,
	}
}

// --------------------------------
// Insert
// --------------------------------
const insertProcessRunLogQuery = `
INSERT INTO process_run_log
(
	process_run_id,
	log_time,
	message,
	type
) VALUES (
	:process_run_id,
	:log_time,
	:message,
	:type
);
SELECT SCOPE_IDENTITY()
`

func (r *ProcessRunLogRepository) Insert(args *models.ProcessRunLog) (int32, error) {
	rows, err := r.DB.Conn.NamedQuery(insertProcessRunLogQuery, args)
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
