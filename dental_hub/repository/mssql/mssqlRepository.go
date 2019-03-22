package mssql

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
