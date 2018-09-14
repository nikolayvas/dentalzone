package repository

import (
	"dental_hub/database"
	m "dental_hub/models"
	mssql "dental_hub/repository/mssql"
	mysql "dental_hub/repository/mysql"

	config "dental_hub/configuration"
)

// Repository is repo
var Repository = repository()

// IRepository is interface that describes db operations
type IRepository interface {

	// dentist
	RegisterDentist(string, string, []byte) (*string, error)
	ActivateDentist(string) error

	Login(string) (*m.Login, error)
	AddPasswordResetConfirmationCode(string, string) error
	ResetPassword([]byte, string, string) error

	GetDentist(*string) (*m.Dentist, error)

	SeedDiagnosis() ([]m.Diagnosis, error)
	SeedManipulations() ([]m.Manipulation, error)
	SeedToothStatuses() ([]m.ToothStatus, error)

	GetPatients(*string) ([]m.Patient, error)
	CreatePatientProfile(m.Patient, *string) error
	UpdatePatientProfile(m.Patient) error
	RemovePatientProfile(*string) error

	GetTeethData(string) (*m.TeethData, error)

	AddToothManipulation(m.ToothAction) error
	RemoveToothManipulation(string) error

	AddToothDiagnosis(m.ToothAction) error
	RemoveToothDiagnosis(string) error

	// patient
	RegisterPatient(string, string, []byte) (*string, error)
	ActivatePatient(string) error

	LoginPatient(string) (*m.Login, error)
	AddPatientPasswordResetConfirmationCode(string, string) error
	ResetPatientPassword([]byte, string, string) error

	// invitation
	InvitePatient(string, string) (*string, error)
	ActivateInvitation(string) error
}

// Repository is repo factory
func repository() IRepository {

	//if else ..
	database.Init()

	if config.GetInstance().DbDriverName == "mssql" {
		return mssql.Repository{Connection: database.DBCon}
	}

	return mysql.Repository{Connection: database.DBCon}
}
