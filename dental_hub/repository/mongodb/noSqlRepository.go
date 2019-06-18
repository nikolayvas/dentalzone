package mongodb

import (
	"context"
	"log"

	minio "github.com/minio/minio-go/v6"
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
	TagsCollection                 string
	TagCollection                  string
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
	TagsCollection:                 "tags",
	TagCollection:                  "tag",
}

// Repository is mongodb implementation of repository
type Repository struct {
	Client      *mongo.Client
	MinioClient *minio.Client
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
