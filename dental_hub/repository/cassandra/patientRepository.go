package cassandra

import (
	m "dental_hub/models"
)

// RegisterPatient registers new patient
func (r *Repository) RegisterPatient(email string, userName string, password []byte) (string, error) {

	return "", nil
}

// ActivatePatient activates alredy registered patient
func (r *Repository) ActivatePatient(id string) error {
	return nil
}

// LoginPatient returns patient details
func (r *Repository) LoginPatient(email string) (*m.Login, error) {
	return nil, nil
}

// AddPatientPasswordResetConfirmationCode inserts new confirmation code in db
func (r *Repository) AddPatientPasswordResetConfirmationCode(email string, code string) error {
	return nil
}

// ResetPatientPassword resets patient password
func (r *Repository) ResetPatientPassword(hashedPassword []byte, email string, code string) error {
	return nil
}
