package cassandra

import (
	ex "dental_hub/exceptions"
	m "dental_hub/models"
	u "dental_hub/utils"
	"encoding/json"
	"sync"
	"time"

	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
)

// Tooth status
type Tooth struct {
	ToothNo          string        `json:"toothno"`
	DiagnosisList    []ToothAction `json:"diagnosislist"`
	ManipulationList []ToothAction `json:"manipulationlist"`
}

// ToothAction could be manipulation or diagnosis
type ToothAction struct {
	RecordID string    `json:"recordid"`
	ActionID int       `json:"actionid"`
	Date     time.Time `json:"date"`
	//IsDeleted   bool
	//DateDeleted time.Time
}

// RegisterDentist registers new user and returns verification code
func (r *Repository) RegisterDentist(email string, userName string, password []byte) (string, error) {

	var id int

	if err := r.Session.Query(`SELECT count(*) FROM dentists WHERE email = ? LIMIT 1`,
		email).Consistency(gocql.One).Scan(&id); err != nil {

		return "", err
	}

	if id > 0 {
		return "", ex.ErrAlreadyExists
	}

	var verificationID uuid.UUID
	verificationID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	verificationIDToString := verificationID.String()

	// Store user information in temporary table and wait for email verification
	if err := r.Session.Query(`INSERT INTO dentistsignup (email, expirationdate, password, username, verificationid) VALUES (?, ?, ?, ?, ?) USING TTL ?`,
		email, time.Now().Add(3*time.Hour), u.EncodingUTF8GetString(password), userName, verificationIDToString, 3*60*60).Exec(); err != nil {
		return "", err
	}

	return verificationIDToString, nil
}

// ActivateDentist activates alredy registered user
func (r *Repository) ActivateDentist(id string) error {

	var userName, email, password string

	if err := r.Session.Query(`SELECT username, email, password FROM dentistsignup WHERE verificationid = ? AND expirationdate > ? LIMIT 1`,
		id, time.Now()).Consistency(gocql.One).Scan(&userName, &email, &password); err != nil {
		switch {
		case err == gocql.ErrNotFound:
			return ex.ErrNotSuch
		default:
			return err
		}
	}

	return r.Session.Query(`INSERT INTO dentists (id, email, password, username, registrationdate) VALUES (uuid(), ?, ?, ?, ?)`,
		email, password, userName, time.Now()).Exec()
}

// Login returns user details
func (r *Repository) Login(email string) (*m.Login, error) {

	var id gocql.UUID
	var password, name string

	if err := r.Session.Query(`SELECT id, password, username FROM dentists WHERE email = ? LIMIT 1`,
		email).Consistency(gocql.One).Scan(&id, &password, &name); err != nil {

		switch {
		case err == gocql.ErrNotFound:
			return nil, ex.ErrNotSuch

		default:
			return nil, err
		}
	}

	login := m.Login{
		ID:       email,
		Email:    email,
		Password: u.EncodingUTF8GetBytes(password),
		Name:     name,
	}

	return &login, nil
}

// AddPasswordResetConfirmationCode insertс new confirmation code in db
func (r *Repository) AddPasswordResetConfirmationCode(email string, code string) error {

	return r.Session.Query(`INSERT INTO dentistresetpassword (dentistemail, code, expirationdate) VALUES (?, ?, ?) USING TTL ?`,
		email, code, time.Now().Add(3*time.Hour), 3*60*60).Exec()
}

// ResetPassword resets user password
func (r *Repository) ResetPassword(hashedPassword []byte, email string, code string) error {

	var _code string

	if err := r.Session.Query(`SELECT code FROM dentistresetpassword WHERE dentistemail = ? LIMIT 1`,
		email).Consistency(gocql.One).Scan(&_code); err != nil {
		switch {
		case err == gocql.ErrNotFound:
			return ex.ErrNotSuch
		default:
			return err
		}
	}

	if _code != code {
		return ex.ErrNotSuch
	} else {
		return r.Session.Query(`UPDATE dentists SET password = ? WHERE email = ?`, u.EncodingUTF8GetString(hashedPassword), email).Exec()
	}
}

// SeedDiagnosis seeds diagnosis
func (r *Repository) SeedDiagnosis() (*[]m.Diagnosis, error) {

	var id, changestatus int
	var diagnosisname string

	diagnosisList := make([]m.Diagnosis, 0)

	iter := r.Session.Query(`SELECT id, changestatus, diagnosisname FROM diagnosis WHERE partitionid=?`, CassandraDbSchema.DiagnosisPartitionKey).Iter()

	for iter.Scan(&id, &changestatus, &diagnosisname) {
		diagnosisList = append(diagnosisList, m.Diagnosis{
			ID:           id,
			ChangeStatus: changestatus,
			Name:         diagnosisname})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return &diagnosisList, nil
}

// SeedManipulations seeds manipulations
func (r *Repository) SeedManipulations() (*[]m.Manipulation, error) {
	var id, changestatus int
	var manipulationname string

	manipulationList := make([]m.Manipulation, 0)

	iter := r.Session.Query(`SELECT id, changestatus, manipulationname FROM manipulations WHERE partitionid=?`, CassandraDbSchema.ManipulationsPartitionKey).Iter()

	for iter.Scan(&id, &changestatus, &manipulationname) {
		manipulationList = append(manipulationList, m.Manipulation{
			ID:           id,
			ChangeStatus: changestatus,
			Name:         manipulationname})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return &manipulationList, nil
}

// SeedToothStatuses seeds tooth statuses
func (r *Repository) SeedToothStatuses() (*[]m.ToothStatus, error) {
	var id int
	var statusname string

	toothStatusList := make([]m.ToothStatus, 0)

	iter := r.Session.Query(`SELECT id, status FROM toothstatus WHERE partitionid=?`, CassandraDbSchema.ToothStatusPartitionKey).Iter()

	for iter.Scan(&id, &statusname) {
		toothStatusList = append(toothStatusList, m.ToothStatus{
			ID:   id,
			Name: statusname})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return &toothStatusList, nil
}

// CreatePatientProfile updates patient
func (r *Repository) CreatePatientProfile(newParient *m.Patient, dentistID string) (string, error) {

	id, err := gocql.RandomUUID()
	if err != nil {
		return "", err
	}

	//todo what if patient with such email already exists?
	if err := r.Session.Query(`INSERT INTO patients (id, email, firstname, middlename, lastname, address, phonenumber, generalinfo, dentists, registrationdate)
		Values(?, ?,?,?,?,?,?,?,?,?)`,
		id,
		newParient.Email,
		newParient.FirstName,
		newParient.MiddleName,
		newParient.LastName,
		newParient.Address,
		newParient.PhoneNumber,
		newParient.GeneralInfo,
		[]string{dentistID},
		time.Now(),
	).Exec(); err != nil {
		return "", err
	}

	if err := r.Session.Query(`INSERT INTO patients_by_id (id, email) Values(?, ?)`,
		id,
		newParient.Email,
	).Exec(); err != nil {
		return "", err
	}

	var patientsList []string

	if err := r.Session.Query(`SELECT patients FROM dentists WHERE email = ? LIMIT 1`,
		dentistID).Consistency(gocql.One).Scan(&patientsList); err != nil {

		return "", err
	}

	patientsList = append(patientsList, newParient.Email)

	err = r.Session.Query(`UPDATE dentists SET patients = ? where email = ?`,
		patientsList,
		dentistID,
	).Exec()

	return id.String(), err
}

// UpdatePatientProfile updates patient
func (r *Repository) UpdatePatientProfile(patient *m.Patient) error {

	return r.Session.Query(`UPDATE patients SET firstname = ?, middlename = ?, lastname = ?, address = ?, phonenumber = ?, generalinfo = ? WHERE email = ?`,
		patient.FirstName, patient.MiddleName, patient.LastName, patient.Address, patient.PhoneNumber, patient.GeneralInfo, patient.Email).Exec()
}

// GetPatients returns patients
func (r *Repository) GetPatients(dentistID string) (*[]m.Patient, error) {
	var patientsIDList []string

	if err := r.Session.Query(`SELECT patients FROM dentists WHERE email = ? LIMIT 1`,
		dentistID).Consistency(gocql.One).Scan(&patientsIDList); err != nil {

		return nil, err
	}

	patients := make([]m.Patient, len(patientsIDList))

	var wg sync.WaitGroup
	wg.Add(len(patientsIDList))

	for i, patientEmail := range patientsIDList {
		go func(i int, email string) {

			var id gocql.UUID

			defer wg.Done()
			if err := r.Session.Query(`SELECT id, email, firstname, middlename, lastname, address, phonenumber, generalinfo, registrationdate from patients WHERE email = ? LIMIT 1`, email).Consistency(gocql.One).Scan(
				&id,
				&(patients[i].Email),
				&(patients[i].FirstName),
				&(patients[i].MiddleName),
				&(patients[i].LastName),
				&(patients[i].Address),
				&(patients[i].PhoneNumber),
				&(patients[i].GeneralInfo),
				&(patients[i].RegistrationDate)); err != nil {
				//todo handle errors
			}

			patients[i].ID = id.String()
		}(i, patientEmail)
	}

	wg.Wait()

	return &patients, nil
}

// RemovePatientProfile removes the patient from the list of patients for the dentist
func (r *Repository) RemovePatientProfile(patientID string, dentistID string) error {
	//todo, we will most probably just unassign instead of delete the patient
	return nil
}

// GetTeethData returns teeth data per patient
func (r *Repository) GetTeethData(patientID string) (*m.TeethData, error) {

	var result = m.TeethData{
		Diagnosis:     []m.ToothAction{},
		Manipulations: []m.ToothAction{},
	}

	patientEmail, err := r.getPatientEmailByID(patientID)

	if err != nil {
		return nil, err
	}

	var teethJSON string

	if err := r.Session.Query(`SELECT teeth_json FROM patients WHERE email = ? LIMIT 1`,
		patientEmail).Consistency(gocql.One).Scan(&teethJSON); err != nil {
		return nil, err
	}

	if teethJSON == "" {
		return &result, nil
	}

	var teeth []Tooth
	if err := json.Unmarshal([]byte(teethJSON), &teeth); err != nil {
		return nil, err
	}

	for _, tooth := range teeth {
		for _, diagnosis := range tooth.DiagnosisList {
			result.Diagnosis = append(result.Diagnosis, m.ToothAction{
				ID:        diagnosis.RecordID,
				ActionID:  diagnosis.ActionID,
				PatientID: patientID,
				ToothNo:   tooth.ToothNo,
				Date:      diagnosis.Date,
			})
		}

		for _, manipulation := range tooth.ManipulationList {
			result.Manipulations = append(result.Manipulations, m.ToothAction{
				ID:        manipulation.RecordID,
				ActionID:  manipulation.ActionID,
				PatientID: patientID,
				ToothNo:   tooth.ToothNo,
				Date:      manipulation.Date,
			})
		}
	}

	return &result, nil
}

// AddToothManipulation adds manipulation
func (r *Repository) AddToothManipulation(manipulation m.ToothAction) error {

	patientEmail, err := r.getPatientEmailByID(manipulation.PatientID)

	if err != nil {
		return err
	}

	var teeth []Tooth
	var teethJSON string

	if err := r.Session.Query(`SELECT teeth_json FROM patients WHERE email = ? LIMIT 1`,
		patientEmail).Consistency(gocql.One).Scan(&teethJSON); err != nil {
		return err
	}

	// JSON to struct conversion
	if teethJSON != "" {
		if err := json.Unmarshal([]byte(teethJSON), teeth); err != nil {
			return err
		}
	}

	_, tooth := FindTooth(&teeth, manipulation.ToothNo)

	if tooth != nil {
		var newList = append(tooth.ManipulationList, ToothAction{
			Date:     time.Now(),
			RecordID: manipulation.ID,
			ActionID: manipulation.ActionID,
		})

		tooth.ManipulationList = newList
	} else {
		tooth = &Tooth{
			ToothNo: manipulation.ToothNo,
			ManipulationList: []ToothAction{
				ToothAction{
					Date:     time.Now(),
					RecordID: manipulation.ID,
					ActionID: manipulation.ActionID,
				},
			},
		}

		teeth = append(teeth, *tooth)
	}

	// struct to JSON conversion
	toothMarshal, _ := json.Marshal(teeth)

	return r.Session.Query(`UPDATE patients SET teeth_json = ? WHERE email = ?`,
		string(toothMarshal), patientEmail).Exec()
}

// RemoveToothManipulation removes manipulation
func (r *Repository) RemoveToothManipulation(manipulation m.ToothAction) error {
	return nil
}

// AddToothDiagnosis adds diagnosis
func (r *Repository) AddToothDiagnosis(diagnosis m.ToothAction) error {
	patientEmail, err := r.getPatientEmailByID(diagnosis.PatientID)

	if err != nil {
		return err
	}

	var teeth []Tooth
	var teethJSON string

	if err := r.Session.Query(`SELECT teeth_json FROM patients WHERE email = ? LIMIT 1`,
		patientEmail).Consistency(gocql.One).Scan(&teethJSON); err != nil {
		return err
	}

	// JSON to struct conversion
	if teethJSON != "" {
		if err := json.Unmarshal([]byte(teethJSON), teeth); err != nil {
			return err
		}
	}

	_, tooth := FindTooth(&teeth, diagnosis.ToothNo)

	if tooth != nil {
		var newList = append(tooth.DiagnosisList, ToothAction{
			Date:     time.Now(),
			RecordID: diagnosis.ID,
			ActionID: diagnosis.ActionID,
		})

		tooth.ManipulationList = newList
	} else {
		tooth = &Tooth{
			ToothNo: diagnosis.ToothNo,
			DiagnosisList: []ToothAction{
				ToothAction{
					Date:     time.Now(),
					RecordID: diagnosis.ID,
					ActionID: diagnosis.ActionID,
				},
			},
		}

		teeth = append(teeth, *tooth)
	}

	// struct to JSON conversion
	toothMarshal, _ := json.Marshal(teeth)

	return r.Session.Query(`UPDATE patients SET teeth_json = ? WHERE email = ?`,
		string(toothMarshal), patientEmail).Exec()
}

// RemoveToothDiagnosis removes diagnosis
func (r *Repository) RemoveToothDiagnosis(diagnosis m.ToothAction) error {
	return nil
}

// InvitePatient resets patient password
func (r *Repository) InvitePatient(dentistID string, patientEmail string) (string, error) {

	return "", nil
}

// ActivateInvitation activates invitation
func (r *Repository) ActivateInvitation(activationID string) error {
	return nil
}

// GetDentist returns dentist data
func (r *Repository) GetDentist(dentistID string) (*m.Dentist, error) {

	var dentist m.Dentist

	if err := r.Session.Query(`SELECT email, username FROM dentists WHERE email = ? LIMIT 1`,
		dentistID).Consistency(gocql.One).Scan(&dentist.Email, &dentist.Name); err != nil {
		switch {
		case err == gocql.ErrNotFound:
			return nil, ex.ErrNotSuch
		default:
			return nil, err
		}
	}

	return &dentist, nil
}

// getPatientEmailByID ...
func (r *Repository) getPatientEmailByID(patientID string) (string, error) {
	var patientEmail string
	id, err := gocql.ParseUUID(patientID)

	if err != nil {
		return "", err
	}

	if err := r.Session.Query(`SELECT email FROM patients_by_id WHERE id = ? LIMIT 1`,
		id).Consistency(gocql.One).Scan(&patientEmail); err != nil {

		return "", err
	}

	return patientEmail, nil
}

// GetAppointments returns appointments per day and dentist
func (r *Repository) GetAppointments(dentistID string, date time.Time) (*[]m.Appointment, error) {

	//return nil, nil

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	justDate := rounded.Format("2006-01-02")
	var appointmentsJSON string

	if err := r.Session.Query(`SELECT appointments FROM schedule WHERE dentistemail = ? and day = ? LIMIT 1`,
		dentistID, justDate).Consistency(gocql.One).Scan(&appointmentsJSON); err != nil {

		switch {
		case err == gocql.ErrNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	var appointments []m.Appointment

	// JSON to struct conversion
	if appointmentsJSON != "" {
		if err := json.Unmarshal([]byte(appointmentsJSON), &appointments); err != nil {
			return nil, err
		}
	}

	return &appointments, nil
}

// UpdateAppointments updates appointments per day and dentist
func (r *Repository) UpdateAppointments(dentistID string, date time.Time, appointments *[]m.Appointment) error {

	// struct to JSON conversion
	appointmentsJSON, _ := json.Marshal(*appointments)

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return r.Session.Query(`UPDATE schedule SET appointments = ? WHERE dentistemail = ? and day = ?`,
		string(appointmentsJSON), dentistID, rounded).Exec()
}

// FindTooth in the collection by toothNo
func FindTooth(teeth *[]Tooth, toothNo string) (int, *Tooth) {
	if teeth == nil {
		return -1, nil
	}

	for i, t := range *teeth {
		if t.ToothNo == toothNo {
			return i, &t
		}
	}

	return -1, nil
}

// FindToothAction in the collection by RecordNo
func FindToothAction(operations []*ToothAction, recordID string) (int, *ToothAction) {
	if operations == nil {
		return -1, nil
	}

	for i, o := range operations {
		if o.RecordID == recordID {
			return i, o
		}
	}

	return -1, nil
}
