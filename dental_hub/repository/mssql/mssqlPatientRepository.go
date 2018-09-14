package mssql

import (
	"database/sql"
	"fmt"

	ex "dental_hub/exceptions"
	m "dental_hub/models"
)

// RegisterPatient registers new patient
func (r Repository) RegisterPatient(email string, userName string, password []byte) (*string, error) {

	var id string
	err := r.Connection.QueryRow("select cast(Id as char(36)) from [dbo].[Patient] where email=?", email).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		//Do nothing, it is expected
	case err != nil:
		return nil, err
	default:
		return nil, ex.ErrAlreadyExists
	}

	rows, err := r.Connection.Query("exec [SignUpPatientRegister] ?, ?, ?", email, userName, password)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var verificationID string

	for rows.Next() {

		err := rows.Scan(&verificationID)

		if err != nil {
			return nil, err
		}
	}

	return &verificationID, nil
}

// ActivatePatient activates alredy registered patient
func (r Repository) ActivatePatient(id string) error {
	res, err := r.Connection.Exec("exec [SignUpPatientActivate] ?", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if count == 0 {
		return ex.ErrNotSuch
	}

	return nil
}

// LoginPatient returns patient details
func (r Repository) LoginPatient(email string) (*m.Login, error) {
	login := m.Login{}

	err := r.Connection.QueryRow("select cast(Id as char(36)), UserName, Email, Password from [dbo].[Patient] where email=?",
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
func (r Repository) AddPatientPasswordResetConfirmationCode(email string, code string) error {
	sql := `delete from [dbo].[ResetPasswordPatient] where PaitentId in (select id from [dbo].[Patient] where [Email] = $1)`

	q, err := r.Connection.Exec(sql, email)

	if err != nil {
		return err
	}

	fmt.Println(q)

	sql = `INSERT INTO [dbo].[ResetPasswordPatient] (PatientId, Code, ExpirationDate) SELECT Id, $1, DATEADD(hour, 3, SYSDATETIME()) from [dbo].[Patient] where [Email] = $2`

	q, err = r.Connection.Exec(sql, code, email)

	if err != nil {
		return err
	}

	return nil
}

// ResetPatientPassword resets patient password
func (r Repository) ResetPatientPassword(hashedPassword []byte, email string, code string) error {

	res, err := r.Connection.Exec("exec [ResetPasswordPatient_SP] ?, ?, ?", hashedPassword, email, code)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if count == 0 {
		return ex.ErrNotSuch
	}

	return nil
}
