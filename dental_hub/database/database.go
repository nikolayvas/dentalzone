package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"

	config "dental_hub/configuration"
)

var (
	// DBCon is the connection handle
	// for the sql database
	DBCon *sql.DB

	// Client is the connection handle
	// for the mongodb
	Client *mongo.Client
)

//Init initialize a database handler that maintains its own pool of idle connections
func Init() {
	if config.GetInstance().DbDriverName == "nosql" {

		client, err := mongo.Connect(context.Background(), config.GetInstance().DbConnectionString)

		if err != nil {
			log.Fatal("Failed to start the Mongo session: " + err.Error())
		}

		err = client.Ping(context.Background(), nil)

		if err != nil {
			log.Fatal("Failed to ping the Mongo server: " + err.Error())
		}

		Client = client
	} else {
		db, err := sql.Open(config.GetInstance().DbDriverName, config.GetInstance().DbConnectionString)

		if err != nil {
			log.Fatal(err)
		}
		DBCon = db
	}
}
