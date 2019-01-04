package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"

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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, config.GetInstance().DbConnectionString)

		defer cancel()
		err = client.Ping(ctx, readpref.Primary())

		if err != nil {
			log.Fatal("Failed to start the Mongo session")
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
