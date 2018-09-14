package models

import "time"

// Dentist model
type Dentist struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Diagnosis model
type Diagnosis struct {
	ID           *string `json:"id"`
	Name         *string `json:"name"`
	ChangeStatus *int    `json:"changeStatus"`
}

// Manipulation model
type Manipulation struct {
	ID           *string `json:"id"`
	Name         *string `json:"name"`
	ChangeStatus *int    `json:"changeStatus"`
}

// ToothStatus model
type ToothStatus struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
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
	ID               *string   `json:"id"`
	FirstName        *string   `json:"firstName"`
	MiddleName       *string   `json:"middleName"`
	LastName         *string   `json:"lastName"`
	Email            *string   `json:"email"`
	Address          *string   `json:"address"`
	PhoneNumber      *string   `json:"phoneNumber"`
	GeneralInfo      *string   `json:"generalInfo"`
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
