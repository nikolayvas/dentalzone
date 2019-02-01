package migrations

import (
	"database/sql"
	db "dental_hub/database"
	"log"
	"strconv"
	"time"

	accessDBwE "github.com/bennof/accessDBwE"

	m "dental_hub/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/satori/go.uuid"
)

var (
	accessDb *sql.DB
)

// Patient read/write model
type Patient struct {
	ID               string    `bson:"_id,omitempty"`
	FirstName        string    `bson:"firstName"`
	MiddleName       string    `bson:"middleName"`
	LastName         string    `bson:"lastName"`
	Email            string    `bson:"email"`
	Address          string    `bson:"address"`
	PhoneNumber      string    `bson:"phoneNumber"`
	GeneralInfo      string    `bson:"generalInfo"`
	RegistrationDate time.Time `bson:"registrationDate"`
	Dentists         []string  `bson:"dentists"`
	Teeth            []*Tooth  `bson:"teeth"`
}

// Tooth ...
type Tooth struct {
	ToothNo       string            `bson:"toothNo"`
	Diagnosis     []*ToothOperation `bson:"diagnosisList"`
	Manipulations []*ToothOperation `bson:"manipulationList"`
}

// ToothOperation could be manipulation or diagnosis
type ToothOperation struct {
	RecordID    string     `bson:"recordID"`
	OperationID int        `bson:"operationID"`
	Date        time.Time  `bson:"date"`
	IsDeleted   *bool      `bson:"isDeleted,omitempty"`
	DateDeleted *time.Time `bson:"dateDeleted,omitempty"`
}

// GetToothByNo ...
func (p Patient) GetToothByNo(toothNo int) (*Tooth, error) {

	if p.Teeth == nil {
		return nil, nil
	}

	for _, tooth := range p.Teeth {
		if tooth.ToothNo == strconv.Itoa(toothNo) {
			return tooth, nil
		}
	}

	return nil, nil
}

// newGuid generates uniqueidentifier
func newGuid() string {
	var recordID uuid.UUID
	recordID, _ = uuid.NewV4()

	return recordID.String()
}

// MigrateAccessDbToSqlServer migrates data from accessDb to SqlServer
func MigrateAccessDbToSqlServer(connectionString string) error {

	err := initDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	/*
		diagnosis, err := seedDiagnosis()
		err = pushDiagnosis(diagnosis)

		if err != nil {
			log.Fatal(err)
		}

		manipulations, err := seedManipulations()
		err = pushManipulations(manipulations)

		if err != nil {
			log.Fatal(err)
		}

		toothstatuses, err := seedToothStatuses()
		err = pushToothStatuses(toothstatuses)

		if err != nil {
			log.Fatal(err)
		}
	*/

	dentistID, err := pushDentist()
	if err != nil {
		log.Fatal(err)
	}

	patients, err := seedPatients()
	if err != nil {
		log.Fatal(err)
	}

	err = pushPatients(dentistID, patients)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// initDbConnection ...
func initDbConnection() error {
	db.Init()

	if accessDb == nil {
		var err error
		accessDb, err = accessDBwE.Open("adodb", "Provider=Microsoft.ACE.OLEDB.12.0;Data Source=C:/Dental.mdb;")
		if err != nil {
			log.Fatal("Failed to start AccessDb session")
		}
	}

	return nil
}

// seedDiagnosis ...
func seedDiagnosis() (*[]m.Diagnosis, error) {
	rows, err := accessDb.Query(`select * from Diagnozi`)
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

// pushDiagnosis ..
func pushDiagnosis(diagnosis *[]m.Diagnosis) error {
	return nil
}

// seedManipulations ...
func seedManipulations() (*[]m.Manipulation, error) {
	rows, err := accessDb.Query(`select * from Manipolacii`)
	if err != nil {
		return nil, err
	}

	manipulationsList := make([]m.Manipulation, 0)

	defer rows.Close()
	for rows.Next() {

		var price float32
		var manipulation m.Manipulation
		err := rows.Scan(
			&manipulation.ID,
			&manipulation.Name,
			&price,
			&manipulation.ChangeStatus)

		if err != nil {
			return nil, err
		}

		manipulationsList = append(manipulationsList, manipulation)
	}

	return &manipulationsList, nil
}

// pushManipulations ...
func pushManipulations(manipulations *[]m.Manipulation) error {
	return nil
}

// seedToothStatuses seeds tooth statuses
func seedToothStatuses() (*[]m.ToothStatus, error) {
	rows, err := accessDb.Query(`select * from Status`)
	if err != nil {
		return nil, err
	}

	toothStatusList := make([]m.ToothStatus, 0)

	defer rows.Close()
	for rows.Next() {

		var toothStatus m.ToothStatus
		err := rows.Scan(
			&toothStatus.ID,
			&toothStatus.Name,
		)

		if err != nil {
			return nil, err
		}

		toothStatusList = append(toothStatusList, toothStatus)
	}

	return &toothStatusList, nil
}

// pushToothStatuses ...
func pushToothStatuses(toothStatuses *[]m.ToothStatus) error {
	return nil
}

// pushDentist ...
func pushDentist() (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("wilber"), 8)
	if err != nil {
		return "", err
	}

	id := newGuid()

	sql := `INSERT INTO Dentist (Id, UserName, Email, Password, RegistrationDate)
			Values($1,$2,$3,$4,$5)`

	_, err = db.DBCon.Exec(sql,
		id,
		"Margarita Tarpanova",
		"a@a.a",
		hashedPassword,
		time.Now(),
	)

	return id, err
}

// seedPatients ..
func seedPatients() (*[]*Patient, error) {

	patients := make([]*Patient, 0)
	m := make(map[int]*Patient)

	rows1, err := accessDb.Query(`SELECT Pacient_Info.ID, Pacient_Info.First_Name, Pacient_Info.Second_Name, Pacient_Info.Last_Name, Pacient_Info.Address, Pacient_Info.Registracia_Data, Pacient_Info.Tel_Domashen, Pacient_Info.Tel_Mobilen, Pacient_Info.Comment, Zaben_Status.Zub_No, Diagnoza_Zub.Diagnoza_Nomer, Diagnoza_Zub.Diagnoza_Data
		FROM (Pacient_Info INNER JOIN Zaben_Status ON Pacient_Info.ID = Zaben_Status.Pacient_ID) INNER JOIN Diagnoza_Zub ON Zaben_Status.Pacient_Zub_ID = Diagnoza_Zub.Pacient_ZubID
		WHERE (((Zaben_Status.StatusID)>1))
		ORDER BY Pacient_Info.ID`)

	if err != nil {
		return nil, err
	}

	defer rows1.Close()
	for rows1.Next() {

		var patientID int
		var firstName sql.NullString
		var secondName sql.NullString
		var lastName sql.NullString
		var address sql.NullString
		var registrationDate time.Time
		var phone1 sql.NullString
		var phone2 sql.NullString
		var comment sql.NullString

		var toothNo int
		var diagnosisNo int
		var diagnosisDate time.Time

		err := rows1.Scan(
			&patientID,
			&firstName,
			&secondName,
			&lastName,
			&address,
			&registrationDate,
			&phone1,
			&phone2,
			&comment,
			&toothNo,
			&diagnosisNo,
			&diagnosisDate,
		)

		if err != nil {
			return nil, err
		}

		patient, ok := m[patientID]

		if !ok {
			patient = &Patient{
				FirstName:        firstName.String,
				MiddleName:       secondName.String,
				LastName:         lastName.String,
				Address:          address.String,
				GeneralInfo:      comment.String,
				RegistrationDate: registrationDate,
			}

			if phone1.String != "" {
				if phone2.String != "" {
					patient.PhoneNumber = phone1.String + ", " + phone2.String
				} else {
					patient.PhoneNumber = phone2.String
				}
			} else {
				if phone2.String != "" {
					patient.PhoneNumber = phone2.String
				}
			}

			patient.Teeth = make([]*Tooth, 0)
			m[patientID] = patient
		}

		tooth, err := patient.GetToothByNo(toothNo)

		if err != nil {
			return nil, err
		}

		if tooth == nil {
			tooth = &Tooth{
				ToothNo:       strconv.Itoa(toothNo),
				Diagnosis:     make([]*ToothOperation, 0),
				Manipulations: make([]*ToothOperation, 0),
			}

			patient.Teeth = append(patient.Teeth, tooth)
		}

		var recordID uuid.UUID
		recordID, err = uuid.NewV4()

		if err != nil {
			return nil, err
		}

		operation := ToothOperation{
			RecordID:    recordID.String(),
			OperationID: diagnosisNo,
			Date:        diagnosisDate,
		}

		tooth.Diagnosis = append(tooth.Diagnosis, &operation)
	}

	rows2, err := accessDb.Query(`SELECT Pacient_Info.ID, Pacient_Info.First_Name, Pacient_Info.Second_Name, Pacient_Info.Last_Name, Pacient_Info.Address, Pacient_Info.Registracia_Data, Pacient_Info.Tel_Domashen, Pacient_Info.Tel_Mobilen, Pacient_Info.Comment, Zaben_Status.Zub_No, Manipolacia_Zub.Manipolacia_Nomer, Manipolacia_Zub.Manipolacia_Data
	FROM (Pacient_Info INNER JOIN Zaben_Status ON Pacient_Info.ID = Zaben_Status.Pacient_ID) INNER JOIN Manipolacia_Zub ON Zaben_Status.Pacient_Zub_ID = Manipolacia_Zub.Pacient_ZubID
	WHERE (((Zaben_Status.StatusID)>1))
	ORDER BY Pacient_Info.ID`)

	if err != nil {
		return nil, err
	}

	defer rows2.Close()

	for rows2.Next() {

		var patientID int
		var firstName sql.NullString
		var secondName sql.NullString
		var lastName sql.NullString
		var address sql.NullString
		var registrationDate time.Time
		var phone1 sql.NullString
		var phone2 sql.NullString
		var comment sql.NullString

		var toothNo int
		var manipilationNo int
		var manipulationDate time.Time

		err := rows2.Scan(
			&patientID,
			&firstName,
			&secondName,
			&lastName,
			&address,
			&registrationDate,
			&phone1,
			&phone2,
			&comment,
			&toothNo,
			&manipilationNo,
			&manipulationDate,
		)

		if err != nil {
			return nil, err
		}

		patient, ok := m[patientID]

		if !ok {
			patient = &Patient{
				FirstName:        firstName.String,
				MiddleName:       secondName.String,
				LastName:         lastName.String,
				Address:          address.String,
				GeneralInfo:      comment.String,
				RegistrationDate: registrationDate,
			}

			if phone1.String != "" {
				if phone2.String != "" {
					patient.PhoneNumber = phone1.String + ", " + phone2.String
				} else {
					patient.PhoneNumber = phone2.String
				}
			} else {
				if phone2.String != "" {
					patient.PhoneNumber = phone2.String
				}
			}

			patient.Teeth = make([]*Tooth, 0)
			m[patientID] = patient
		}

		tooth, err := patient.GetToothByNo(toothNo)

		if err != nil {
			return nil, err
		}

		if tooth == nil {
			tooth = &Tooth{
				ToothNo:       strconv.Itoa(toothNo),
				Diagnosis:     make([]*ToothOperation, 0),
				Manipulations: make([]*ToothOperation, 0),
			}

			patient.Teeth = append(patient.Teeth, tooth)
		}

		var recordID uuid.UUID
		recordID, err = uuid.NewV4()

		if err != nil {
			return nil, err
		}

		operation := ToothOperation{
			RecordID:    recordID.String(),
			OperationID: manipilationNo,
			Date:        manipulationDate,
		}

		tooth.Manipulations = append(tooth.Manipulations, &operation)
	}

	for _, v := range m {
		patients = append(patients, v)
	}

	return &patients, nil
}

// PushPatients ...
func pushPatients(dentistID string, patients *[]*Patient) error {

	for _, newParient := range *patients {

		id := newGuid()

		sql := `INSERT INTO PatientInfo (Id, FirstName, MiddleName, LastName, Email, Address, PhoneNumber, GeneralInfo, DentistId)
		Values($1,$2,$3,$4,$5,$6,$7,$8,$9)`

		_, err := db.DBCon.Exec(sql,
			id,
			newParient.FirstName,
			newParient.MiddleName,
			newParient.LastName,
			newParient.Email,
			newParient.Address,
			newParient.PhoneNumber,
			newParient.GeneralInfo,
			dentistID)

		if err != nil {
			return err
		}

		for _, tooth := range newParient.Teeth {
			toothID := newGuid()

			sql := `INSERT INTO ToothCurrentStatus (ToothID, ToothNo, PatientId)
			Values($1,$2,$3)`

			_, err = db.DBCon.Exec(sql,
				toothID,
				tooth.ToothNo,
				id)

			if err != nil {
				return err
			}

			for _, diagnosis := range tooth.Diagnosis {
				diagnosisID := newGuid()

				sql := `INSERT INTO ToothDiagnosis (Id, DiagnosisId, Date, ToothId)
					Values($1,$2,$3,$4)`

				_, err = db.DBCon.Exec(sql,
					diagnosisID,
					diagnosis.OperationID,
					diagnosis.Date,
					toothID)

				if err != nil {
					return err
				}
			}

			for _, manipulation := range tooth.Manipulations {
				manipulationID := newGuid()

				sql := `INSERT INTO ToothManipulation (Id, ManipulationId, Date, ToothId)
				Values($1,$2,$3,$4)`

				_, err = db.DBCon.Exec(sql,
					manipulationID,
					manipulation.OperationID,
					manipulation.Date,
					toothID)

				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}
