package mysql

import (
	"database/sql"
	"log"
)

// Repository is sql server implementation of repository
type Repository struct {
	Connection *sql.DB
}

// Init connction
func (r *Repository) Init(driverName string, connectionString string) {
	db, err := sql.Open(driverName, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	r.Connection = db
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
