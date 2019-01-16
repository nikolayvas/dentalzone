package models

import (
	"time"
)

// Dentist model
type Dentist struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password []byte `json:"password"`
}

// Diagnosis model
type Diagnosis struct {
	ID           int    `json:"id"`
	Name         string `json:"name" bson:"diagnosisName"`
	ChangeStatus int    `json:"changeStatus"`
}

// Manipulation model
type Manipulation struct {
	ID           int    `json:"id"`
	Name         string `json:"name" bson:"manipulationName"`
	ChangeStatus int    `json:"changeStatus"`
}

// ToothStatus model
type ToothStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name" bson:"status"`
}

// Login model
type Login struct {
	ID       string
	Email    string
	Name     string
	Password []byte
}

// Patient model
type Patient struct {
	ID               string    `json:"id"`
	FirstName        string    `json:"firstName"`
	MiddleName       string    `json:"middleName"`
	LastName         string    `json:"lastName"`
	Email            string    `json:"email"`
	Address          string    `json:"address"`
	PhoneNumber      string    `json:"phoneNumber"`
	GeneralInfo      string    `json:"generalInfo"`
	RegistrationDate time.Time `json:"registrationDate"`
}

// ToothAction model
type ToothAction struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patientId"`
	ToothNo   string    `json:"toothNo"`
	ActionID  int       `json:"actionId"`
	Date      time.Time `json:"date"`
}

// TeethData model
type TeethData struct {
	Diagnosis     []ToothAction `json:"diagnosisList"`
	Manipulations []ToothAction `json:"manipulationList"`
}

// Appointment model
type Appointment struct {
	Date      time.Time `json:"date"`
	PatientID string    `json:"patientID"`
}

// Appointments model
type Appointments struct {
	Day          time.Time     `json:"day"`
	Appointments []Appointment `json:"appointments"`
}
