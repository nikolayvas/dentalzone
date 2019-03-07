package cassandra

import (
	"github.com/gocql/gocql"
)

// Repository is cassandra implementation of repository
type Repository struct {
	Session *gocql.Session
}
