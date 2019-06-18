package mssql

import (
	"database/sql"
	"log"

	minio "github.com/minio/minio-go/v6"
)

// Repository is sql server implementation of repository
type Repository struct {
	Connection  *sql.DB
	MinioClient *minio.Client
}

// Init connction
func (r *Repository) Init(driverName string, connectionString string) {
	db, err := sql.Open(driverName, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	r.Connection = db

	//minio setup
	endpoint := "localhost:9000"
	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	r.MinioClient = minioClient
}
