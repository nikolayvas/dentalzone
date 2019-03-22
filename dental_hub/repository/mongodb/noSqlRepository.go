package mongodb

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// DbSchema stores db's objects names
type DbSchema struct {
	DatabaseName                   string
	DentistSignUpCollection        string
	DentistCollection              string
	DentistResetPasswordCollection string
	DiagnosisCollection            string
	ManipulationsCollection        string
	ToothStatusesCollection        string
	PatientCollection              string
	ScheduleCollection             string
}

// MongoDbSchema specifies db schema
var MongoDbSchema = DbSchema{
	DatabaseName:                   "dental_hub",
	DentistSignUpCollection:        "dentistSignUp",
	DentistCollection:              "dentist",
	DentistResetPasswordCollection: "dentistResetPassword",
	DiagnosisCollection:            "diagnosis",
	ManipulationsCollection:        "manipulations",
	ToothStatusesCollection:        "toothStatuses",
	PatientCollection:              "patient",
	ScheduleCollection:             "schedule",
}

// Repository is mongodb implementation of repository
type Repository struct {
	Client *mongo.Client
}

// Init connction
func (r *Repository) Init(connectionString string) {
	client, err := mongo.Connect(context.Background(), connectionString)

	if err != nil {
		log.Fatal("Failed to start the Mongo session: " + err.Error())
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("Failed to ping the Mongo server: " + err.Error())
	}

	r.Client = client
}
