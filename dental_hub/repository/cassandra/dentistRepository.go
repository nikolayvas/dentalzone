package cassandra

import (
	ex "dental_hub/exceptions"
	m "dental_hub/models"
	u "dental_hub/utils"
	"time"

	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
)

// RegisterDentist registers new user and returns verification code
func (r *Repository) RegisterDentist(email string, userName string, password []byte) (string, error) {

	var id gocql.UUID

	if err := r.Session.Query(`SELECT id FROM dentists WHERE email = ? LIMIT 1`,
		email).Consistency(gocql.One).Scan(&id); err != nil {
		switch {
		case err == gocql.ErrNotFound:
			// Do nothing, the email is available for next registration
		case err != nil:
			return "", err
		default:
			return "", ex.ErrAlreadyExists
		}
	}

	// if already have unprocessed requests for registration with such email, remove them
	if err := r.Session.Query(`DELETE FROM dentistsignup WHERE email = ?`, email).Exec(); err != nil {
		return "", err
	}

	var verificationID uuid.UUID
	verificationID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	verificationIDToString := verificationID.String()

	// Store user information in temporary table and wait for email verification
	if err := r.Session.Query(`INSERT INTO dentistsignup (email, expirationdate, password, username, verificationid) VALUES (?, ?, ?, ?, ?)`,
		email, time.Now().Add(3*time.Hour), u.EncodingUTF8GetString(password), userName, verificationIDToString).Exec(); err != nil {
		return "", err
	}

	return verificationIDToString, nil
}

// ActivateDentist activates alredy registered user
func (r *Repository) ActivateDentist(id string) error {
	return nil
}

// Login returns user details
func (r *Repository) Login(email string) (*m.Login, error) {
	return nil, nil
}

// AddPasswordResetConfirmationCode insertс new confirmation code in db
func (r *Repository) AddPasswordResetConfirmationCode(email string, code string) error {
	return nil
}

// ResetPassword resets user password
func (r *Repository) ResetPassword(hashedPassword []byte, email string, code string) error {
	return nil
}

// SeedDiagnosis seeds diagnosis
func (r *Repository) SeedDiagnosis() (*[]m.Diagnosis, error) {
	return nil, nil
}

// SeedManipulations seeds manipulations
func (r *Repository) SeedManipulations() (*[]m.Manipulation, error) {
	return nil, nil
}

// SeedToothStatuses seeds tooth statuses
func (r *Repository) SeedToothStatuses() (*[]m.ToothStatus, error) {
	return nil, nil
}

// CreatePatientProfile updates patient
func (r *Repository) CreatePatientProfile(newParient m.Patient, dentistID string) (string, error) {
	return "", nil
}

// UpdatePatientProfile updates patient
func (r *Repository) UpdatePatientProfile(patient m.Patient) error {
	return nil
}

// GetPatients returns patients
func (r *Repository) GetPatients(dentistID string) (*[]m.Patient, error) {
	return nil, nil
}

// RemovePatientProfile removes the patient from the list of patients for the dentist
func (r *Repository) RemovePatientProfile(patientID string, dentistID string) error {
	return nil
}

// GetTeethData returns teeth data per patient
func (r *Repository) GetTeethData(patientID string) (*m.TeethData, error) {
	return nil, nil
}

// AddToothManipulation adds manipulation
func (r *Repository) AddToothManipulation(manipulation m.ToothAction) error {
	return nil
}

// RemoveToothManipulation removes manipulation
func (r *Repository) RemoveToothManipulation(manipulation m.ToothAction) error {
	return nil
}

// AddToothDiagnosis adds diagnosis
func (r *Repository) AddToothDiagnosis(diagnosis m.ToothAction) error {
	return nil
}

// RemoveToothDiagnosis removes diagnosis
func (r *Repository) RemoveToothDiagnosis(diagnosis m.ToothAction) error {
	return nil
}

// InvitePatient resets patient password
func (r *Repository) InvitePatient(dentistID string, patientEmail string) (string, error) {

	return "", nil

	//TODO what if already assigned
}

// ActivateInvitation activates invitation
func (r *Repository) ActivateInvitation(activationID string) error {
	return nil
}

// GetDentist returns dentist data
func (r *Repository) GetDentist(dentistID string) (*m.Dentist, error) {
	return nil, nil
}

// GetAppointments returns appointments per day and dentist
func (r *Repository) GetAppointments(dentistID string, date time.Time) (*[]m.Appointment, error) {
	return nil, nil
}

// UpdateAppointments updates appointments per day and dentist
func (r *Repository) UpdateAppointments(dentistID string, date time.Time, appointments *[]m.Appointment) error {
	return nil
}
