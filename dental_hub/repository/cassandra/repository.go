package cassandra

import (
	"log"

	"github.com/gocql/gocql"
	minio "github.com/minio/minio-go/v6"
)

// DbSchema stores db's objects names
type DbSchema struct {
	DatabaseName              string
	DiagnosisPartitionKey     string
	ManipulationsPartitionKey string
	ToothStatusPartitionKey   string
}

// CassandraDbSchema specifies db schema
var CassandraDbSchema = DbSchema{
	DatabaseName:              "dental_hub",
	DiagnosisPartitionKey:     "diagnosisKey",
	ManipulationsPartitionKey: "manipulationsKey",
	ToothStatusPartitionKey:   "toothStatusKey",
}

// Repository is cassandra implementation of repository
type Repository struct {
	Session     *gocql.Session
	MinioClient *minio.Client
}

// Init connction
func (r *Repository) Init() {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "dental_db"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to create cassandra session: " + err.Error())
	}

	r.Session = session

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

	log.Printf("%#v\n", minioClient) // minioClient is now setup
}
