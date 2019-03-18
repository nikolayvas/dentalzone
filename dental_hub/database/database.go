package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/gocql/gocql"
	"github.com/mongodb/mongo-go-driver/mongo"

	config "dental_hub/configuration"
)

var (
	// DBCon is the connection handle for the sql database
	DBCon *sql.DB

	// Client is the connection handle for the mongodb
	Client *mongo.Client

	// Session is cassandra db session
	Session *gocql.Session
)

//Init initialize a database handler that maintains its own pool of idle connections
func Init() {
	if config.GetInstance().DbDriverName == "mongoDb" {

		client, err := mongo.Connect(context.Background(), config.GetInstance().DbConnectionString)

		if err != nil {
			log.Fatal("Failed to start the Mongo session: " + err.Error())
		}

		err = client.Ping(context.Background(), nil)

		if err != nil {
			log.Fatal("Failed to ping the Mongo server: " + err.Error())
		}

		Client = client
	} else if config.GetInstance().DbDriverName == "cassandra" {
		cluster := gocql.NewCluster("localhost")
		cluster.Keyspace = "dental_db"
		cluster.Consistency = gocql.Quorum
		session, err := cluster.CreateSession()
		if err != nil {
			log.Fatal("Failed to create cassandra session: " + err.Error())
		}

		Session = session
	} else {
		db, err := sql.Open(config.GetInstance().DbDriverName, config.GetInstance().DbConnectionString)

		if err != nil {
			log.Fatal(err)
		}
		DBCon = db
	}
}
