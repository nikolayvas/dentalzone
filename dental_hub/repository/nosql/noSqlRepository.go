package nosql

import (
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
