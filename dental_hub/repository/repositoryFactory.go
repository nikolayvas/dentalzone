package repository

import (
	m "dental_hub/models"
	"dental_hub/repository/cassandra"
	"dental_hub/repository/mongodb"
	"dental_hub/repository/mssql"
	"dental_hub/repository/mysql"
	"io"
	"log"
	"time"

	config "dental_hub/configuration"
)

// Repository is repo
var Repository = repository()

// IRepository is interface that describes db operations
type IRepository interface {

	// ------------------dentist----------------------
	RegisterDentist(string, string, []byte) (string, error)
	ActivateDentist(string) error

	Login(string) (*m.Login, error)
	AddPasswordResetConfirmationCode(string, string) error
	ResetPassword([]byte, string, string) error

	GetDentist(string) (*m.Dentist, error)

	SeedDiagnosis() (*[]m.Diagnosis, error)
	SeedManipulations() (*[]m.Manipulation, error)
	SeedToothStatuses() (*[]m.ToothStatus, error)

	GetPatients(string) (*[]m.Patient, error)
	CreatePatientProfile(*m.Patient, string) (string, error)
	UpdatePatientProfile(*m.Patient) error
	RemovePatientProfile(string, string) error

	GetTeethData(string) (*m.TeethData, error)

	AddToothManipulation(m.ToothAction) error
	RemoveToothManipulation(m.ToothAction) error

	AddToothDiagnosis(m.ToothAction) error
	RemoveToothDiagnosis(m.ToothAction) error

	GetAppointments(string, time.Time) (*[]m.Appointment, error)
	UpdateAppointments(string, time.Time, *[]m.Appointment) error

	// ----------------patient---------------------
	RegisterPatient(string, string, []byte) (string, error)
	ActivatePatient(string) error

	LoginPatient(string) (*m.Login, error)
	AddPatientPasswordResetConfirmationCode(string, string) error
	ResetPatientPassword([]byte, string, string) error

	// invitation
	InvitePatient(string, string) (string, error)
	ActivateInvitation(string) error

	// ----------------imaging--------------------

	// insert image and specify the tags
	InsertImage(string /*patient*/, io.Reader /*imageFileLocation*/, []string /*tags*/, int64 /*file size*/) error

	// get image id's by tags
	GetImageIdsByTags(string /*patient*/, []string /*tags*/) ([]string /*S3/minio file id's*/, error)

	// get S3/minio image by id
	GetImage(string /*patient*/, string /*S3/minio image id*/) (io.Reader, error)
}

// Repository is repo factory
func repository() IRepository {

	conf := config.GetInstance()

	if conf.DbDriverName == "mssql" {
		repo := mssql.Repository{}
		repo.Init(conf.DbDriverName, conf.DbConnectionString)
		return &repo
	} else if conf.DbDriverName == "mysql" {
		repo := mysql.Repository{}
		repo.Init(conf.DbDriverName, conf.DbConnectionString)
		return &repo
	} else if conf.DbDriverName == "cassandra" {
		repo := cassandra.Repository{}
		repo.Init()
		return &repo
	} else if conf.DbDriverName == "mongoDb" {
		repo := mongodb.Repository{}
		repo.Init(conf.DbConnectionString)
		return &repo
	} else {
		log.Fatal("Unsuppoerted driver!")
		return nil
	}
}
