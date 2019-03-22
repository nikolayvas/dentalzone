package cassandra

import (
	"log"

	"github.com/gocql/gocql"
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
	Session *gocql.Session
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
}
