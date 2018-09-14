package mssql

import (
	"database/sql"
)

// Repository is sql server implementation of repository
type Repository struct {
	Connection *sql.DB
}
