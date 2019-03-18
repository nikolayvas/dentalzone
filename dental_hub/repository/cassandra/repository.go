package cassandra

import (
	"github.com/gocql/gocql"
)

// DbSchema stores db's objects names
type DbSchema struct {
	DatabaseName              string
	DiagnosisPartitionKey     string
	ManipulationsPartitionKey string
	ToothStatusPartitionKey   string
}

// MongoDbSchema specifies db schema
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
