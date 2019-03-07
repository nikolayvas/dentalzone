package mysql

import (
	"database/sql"

	ex "dental_hub/exceptions"
	m "dental_hub/models"
)

// RegisterPatient registers new patient
func (r *Repository) RegisterPatient(email string, userName string, password []byte) (string, error) {

	var id string
	err := r.Connection.QueryRow("select Id from Patient where email=?", email).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		//Do nothing, it is expected
	case err != nil:
		return "", err
	default:
		return "", ex.ErrAlreadyExists
	}

	rows, err := r.Connection.Query("call signup_patient_register(?, ?, ?)", email, userName, password)

	if err != nil {
		return "", err
	}

	defer rows.Close()
	var verificationID string

	for rows.Next() {

		err := rows.Scan(&verificationID)

		if err != nil {
			return "", err
		}
	}

	return verificationID, nil
}

// ActivatePatient activates alredy registered patient
func (r *Repository) ActivatePatient(id string) error {
	rows, err := r.Connection.Query("call signup_patient_activate(?)", id)
	if err != nil {
		return err
	}

	defer rows.Close()

	count, err := GetSPAffectedRows(rows)
	if count == 0 {
		return ex.ErrNotSuch
	}

	return nil
}

// LoginPatient returns patient details
func (r *Repository) LoginPatient(email string) (*m.Login, error) {
	login := m.Login{}

	err := r.Connection.QueryRow("select Id, UserName, Email, Password from Patient where email=?",
		email).Scan(&login.ID, &login.Name, &login.Email, &login.Password)

	switch {
	case err == sql.ErrNoRows:
		return nil, ex.ErrNotSuch
	case err != nil:
		return nil, err
	default:
		return &login, nil
	}
}

// AddPatientPasswordResetConfirmationCode inserts new confirmation code in db
func (r *Repository) AddPatientPasswordResetConfirmationCode(email string, code string) error {

	_, err := r.Connection.Exec("call add_patient_password_reset_confirmation_code(?, ?)", email, code)
	if err != nil {
		return err
	}

	return nil
}

// ResetPatientPassword resets patient password
func (r *Repository) ResetPatientPassword(hashedPassword []byte, email string, code string) error {

	rows, err := r.Connection.Query("call reset_password_patient_sp(?, ?, ?)", hashedPassword, email, code)
	if err != nil {
		return err
	}

	defer rows.Close()

	count, err := GetSPAffectedRows(rows)
	if count == 0 {
		return ex.ErrNotSuch
	}

	return nil
}
