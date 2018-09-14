package database

import (
	"database/sql"
	"log"

	config "dental_hub/configuration"
)

var (
	// DBCon is the connection handle
	// for the database
	DBCon *sql.DB
)

//Init initialize a database handler that maintains its own pool of idle connections
func Init() {
	db, err := sql.Open(config.GetInstance().DbDriverName, config.GetInstance().DbConnectionString)
	//db, err := sql.Open("mssql", "server=localhost;user id=sa;password=Manager123;database=Dental;connection timeout=30;")
	//db, err := sql.Open("mysql", "root:root@/dental?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	DBCon = db
}
