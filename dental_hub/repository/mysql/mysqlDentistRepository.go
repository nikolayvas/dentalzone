package mysql

import (
	"database/sql"
	"errors"
	"time"

	"dental_hub/database"
	ex "dental_hub/exceptions"
	m "dental_hub/models"
)

// RegisterDentist registers new user
func (r Repository) RegisterDentist(email string, userName string, password []byte) (string, error) {

	var id string
	err := r.Connection.QueryRow("select Id from Dentist where email=?", email).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		//Do nothing, it is expected
	case err != nil:
		return "", err
	default:
		return "", ex.ErrAlreadyExists
	}

	rows, err := r.Connection.Query("call signup_register(?, ?, ?)", email, userName, password)

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

// ActivateDentist activates alredy registered user
func (r Repository) ActivateDentist(id string) error {
	rows, err := r.Connection.Query("call signup_activate(?)", id)
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

// Login returns user details
func (r Repository) Login(email string) (*m.Login, error) {
	login := m.Login{}

	err := r.Connection.QueryRow("select Id, UserName, Email, Password from Dentist where email=?",
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

// AddPasswordResetConfirmationCode insertс new confirmation code in db
func (r Repository) AddPasswordResetConfirmationCode(email string, code string) error {

	_, err := r.Connection.Exec("call add_password_reset_confirmation_code(?, ?)", email, code)

	return err
}

// ResetPassword resets user password
func (r Repository) ResetPassword(hashedPassword []byte, email string, code string) error {

	rows, err := r.Connection.Query("call reset_password_sp(?, ?, ?)", hashedPassword, email, code)
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

// SeedDiagnosis seeds diagnosis
func (r Repository) SeedDiagnosis() (*[]m.Diagnosis, error) {
	rows, err := r.Connection.Query(`select * from Diagnosis`)
	if err != nil {
		return nil, err
	}

	diagnosisList := make([]m.Diagnosis, 0)

	defer rows.Close()
	for rows.Next() {

		var diagnosis m.Diagnosis
		err := rows.Scan(
			&diagnosis.ID,
			&diagnosis.Name,
			&diagnosis.ChangeStatus)

		if err != nil {
			return nil, err
		}

		diagnosisList = append(diagnosisList, diagnosis)
	}

	return &diagnosisList, nil
}

// SeedManipulations seeds manipulations
func (r Repository) SeedManipulations() (*[]m.Manipulation, error) {

	rows, err := r.Connection.Query(`select * from Manipulations`)
	if err != nil {
		return nil, err
	}

	manipulationList := make([]m.Manipulation, 0)

	defer rows.Close()
	for rows.Next() {

		var manipilation m.Manipulation
		err := rows.Scan(
			&manipilation.ID,
			&manipilation.Name,
			&manipilation.ChangeStatus)

		if err != nil {
			return nil, err
		}

		manipulationList = append(manipulationList, manipilation)
	}

	return &manipulationList, nil
}

// SeedToothStatuses seeds tooth statuses
func (r Repository) SeedToothStatuses() (*[]m.ToothStatus, error) {

	rows, err := r.Connection.Query(`select * from ToothStatus`)
	if err != nil {
		return nil, err
	}

	statusesList := make([]m.ToothStatus, 0)

	defer rows.Close()
	for rows.Next() {

		var status m.ToothStatus
		err := rows.Scan(
			&status.ID,
			&status.Name)

		if err != nil {
			return nil, err
		}

		statusesList = append(statusesList, status)
	}

	return &statusesList, nil
}

// GetPatients returns patients
func (r Repository) GetPatients(dentistID string) (*[]m.Patient, error) {

	patients := make([]m.Patient, 0)

	rows, err := database.DBCon.Query(
		`select 
			Id,
			FirstName,
			MiddleName,
			LastName,
			Email,
			Address,
			PhoneNumber,
			GeneralInfo,
			RegistrationDate
		 from patientinfo where dentistId = ? and (IsDeleted = 0 OR IsDeleted is NULL)`, dentistID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var patient m.Patient
		err := rows.Scan(
			&patient.ID,
			&patient.FirstName,
			&patient.MiddleName,
			&patient.LastName,
			&patient.Email,
			&patient.Address,
			&patient.PhoneNumber,
			&patient.GeneralInfo,
			&patient.RegistrationDate)

		if err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	err = rows.Err()
	switch {
	case err != nil:
		return nil, err
	default:
	}

	return &patients, nil
}

// UpdatePatientProfile updates patient
func (r Repository) UpdatePatientProfile(patient m.Patient) error {
	sql := `UPDATE PatientInfo SET 
				FirstName = ?,
				MiddleName = ?,
				LastName = ?,
				Email = ?,
				Address = ?,
				PhoneNumber = ?,
				GeneralInfo = ?
			WHERE Id = ?`

	_, err := database.DBCon.Exec(sql,
		patient.FirstName,
		patient.MiddleName,
		patient.LastName,
		patient.Email,
		patient.Address,
		patient.PhoneNumber,
		patient.GeneralInfo,
		patient.ID)

	return err
}

// CreatePatientProfile creates patient
func (r Repository) CreatePatientProfile(newParient m.Patient, dentistID string) (string, error) {
	sql := `INSERT INTO PatientInfo (Id, FirstName, MiddleName, LastName, Email, Address, PhoneNumber, GeneralInfo, DentistId)
			Values(?,?,?,?,?,?,?,?,?)`

	_, err := database.DBCon.Exec(sql,
		newParient.ID,
		newParient.FirstName,
		newParient.MiddleName,
		newParient.LastName,
		newParient.Email,
		newParient.Address,
		newParient.PhoneNumber,
		newParient.GeneralInfo,
		dentistID)

	return newParient.ID, err
}

// RemovePatientProfile updates patient
func (r Repository) RemovePatientProfile(patientID string, dentistID string) error {
	sql := `UPDATE PatientInfo SET 
				IsDeleted = 1
			WHERE Id= ?`

	_, err := database.DBCon.Exec(sql, patientID)

	return err
}

// GetTeethData returns teeth data per patient
func (r Repository) GetTeethData(patientID string) (*m.TeethData, error) {
	diagnosiesList := make([]m.ToothAction, 0)
	manipulationsList := make([]m.ToothAction, 0)

	rows, err := database.DBCon.Query("call get_teethData(?)", patientID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var diagnosis m.ToothAction
		err := rows.Scan(
			&diagnosis.ID,
			&diagnosis.ActionID,
			&diagnosis.Date,
			&diagnosis.ToothNo)

		if err != nil {
			return nil, err
		}

		diagnosiesList = append(diagnosiesList, diagnosis)
	}

	if !rows.NextResultSet() {
		return nil, errors.New("Missing result set data")
	}

	for rows.Next() {

		var manipulation m.ToothAction
		err := rows.Scan(
			&manipulation.ID,
			&manipulation.ActionID,
			&manipulation.Date,
			&manipulation.ToothNo)

		if err != nil {
			return nil, err
		}

		manipulationsList = append(manipulationsList, manipulation)
	}

	result := m.TeethData{Diagnosis: diagnosiesList, Manipulations: manipulationsList}

	err = rows.Err()
	switch {
	case err != nil:
		return nil, err
	default:
		return &result, nil
	}
}

// AddToothManipulation adds manipulation
func (r Repository) AddToothManipulation(manipulation m.ToothAction) error {
	_, err := r.Connection.Exec("call add_manupulation(?, ?, ?, ?)",
		manipulation.ID,
		manipulation.PatientID,
		manipulation.ToothNo,
		manipulation.ActionID)

	return err
}

// RemoveToothManipulation removes manipulation
func (r Repository) RemoveToothManipulation(manipulation m.ToothAction) error {
	_, err := r.Connection.Exec("call remove_tooth_manipulation(?)", manipulation.ID)

	return err
}

// AddToothDiagnosis adds diagnosis
func (r Repository) AddToothDiagnosis(diagnosis m.ToothAction) error {
	_, err := r.Connection.Exec("call add_diagnosis(?, ?, ?, ?)",
		diagnosis.ID,
		diagnosis.PatientID,
		diagnosis.ToothNo,
		diagnosis.ActionID)

	return err
}

// RemoveToothDiagnosis removes diagnosis
func (r Repository) RemoveToothDiagnosis(diagnosis m.ToothAction) error {
	_, err := r.Connection.Exec("call remove_tooth_diagnosis(?)", diagnosis.ID)

	return err
}

// InvitePatient resets patient password
func (r Repository) InvitePatient(dentistID string, patientEmail string) (string, error) {

	rows, err := r.Connection.Query("call invite_patient(?, ?)", dentistID, patientEmail)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	var invitationID string

	for rows.Next() {
		err := rows.Scan(
			&invitationID)

		if err != nil {
			return "", err
		}
	}

	if invitationID == "" {
		return "", ex.ErrNotSuch
	}

	return invitationID, nil

	//TODO what if already assigned
}

// ActivateInvitation activates invitation
func (r Repository) ActivateInvitation(activationID string) error {
	rows, err := r.Connection.Query("call activate_invitation(?)", activationID)
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

// GetDentist returns dentist data
func (r Repository) GetDentist(id string) (*m.Dentist, error) {
	var dentist m.Dentist

	err := r.Connection.QueryRow("select UserName, Email from Dentist where id=?",
		id).Scan(&dentist.Name, &dentist.Email)

	return &dentist, err
}

// GetAppointments returns appointments for day
func (r Repository) GetAppointments(patientID string, day time.Time) (*[]m.Appointment, error) {
	return nil, nil
}

// UpdateAppointments updates appointments for day
func (r Repository) UpdateAppointments(patientID string, day time.Time, appointments *[]m.Appointment) error {
	return nil
}
