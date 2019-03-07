package repository

import (
	"dental_hub/database"
	m "dental_hub/models"
	"dental_hub/repository/mongodb"
	"dental_hub/repository/mssql"
	"dental_hub/repository/mysql"
	"time"

	config "dental_hub/configuration"
)

// Repository is repo
var Repository = repository()

// IRepository is interface that describes db operations
type IRepository interface {

	// dentist
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
	CreatePatientProfile(m.Patient, string) (string, error)
	UpdatePatientProfile(m.Patient) error
	RemovePatientProfile(string, string) error

	GetTeethData(string) (*m.TeethData, error)

	AddToothManipulation(m.ToothAction) error
	RemoveToothManipulation(m.ToothAction) error

	AddToothDiagnosis(m.ToothAction) error
	RemoveToothDiagnosis(m.ToothAction) error

	GetAppointments(string, time.Time) (*[]m.Appointment, error)
	UpdateAppointments(string, time.Time, *[]m.Appointment) error

	// patient
	RegisterPatient(string, string, []byte) (string, error)
	ActivatePatient(string) error

	LoginPatient(string) (*m.Login, error)
	AddPatientPasswordResetConfirmationCode(string, string) error
	ResetPatientPassword([]byte, string, string) error

	// invitation
	InvitePatient(string, string) (string, error)
	ActivateInvitation(string) error
}

// Repository is repo factory
func repository() IRepository {

	//if else ..
	database.Init()

	if config.GetInstance().DbDriverName == "mssql" {
		return &mssql.Repository{Connection: database.DBCon}
	} else if config.GetInstance().DbDriverName == "mysql" {
		return &mysql.Repository{Connection: database.DBCon}
	}

	return &mongodb.Repository{Client: database.Client}
}
