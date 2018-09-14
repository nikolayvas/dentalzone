package mysql

import (
	"database/sql"
)

// Repository is sql server implementation of repository
type Repository struct {
	Connection *sql.DB
}

// GetSPAffectedRows returns affected Rows returned from the result
func GetSPAffectedRows(rows *sql.Rows) (int, error) {
	var affctedRows int

	for rows.Next() {

		err := rows.Scan(&affctedRows)

		if err != nil {
			return 0, err
		}
	}

	return affctedRows, nil
}
