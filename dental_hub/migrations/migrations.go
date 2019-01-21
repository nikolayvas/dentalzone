package migrations

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	accessDBwE "github.com/bennof/accessDBwE"
	"gopkg.in/mgo.v2/bson"

	m "dental_hub/models"

	config "dental_hub/configuration"
	nosql "dental_hub/repository/nosql"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"

	"golang.org/x/crypto/bcrypt"

	"github.com/satori/go.uuid"
)

var (
	accessDb *sql.DB
	mongoDb  *mongo.Client
)

// MigrateAccessDb migrates data from accessDb to MongoDb
func MigrateAccessDb(connectionString string) error {

	_ = initDbConnection()

	/*
		var diagnosis *[]m.Diagnosis
		diagnosis, err := seedDiagnosis()
		err = pushDiagnosis(diagnosis)

		var manipulations *[]m.Manipulation
		manipulations, err := seedManipulations()
		err = pushManipulations(manipulations)

		var toothstatuses *[]m.ToothStatus
		toothstatuses, err := seedToothStatuses()
		err = pushToothStatuses(toothstatuses)
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

	if accessDb == nil {
		var err error
		accessDb, err = accessDBwE.Open("adodb", "Provider=Microsoft.ACE.OLEDB.12.0;Data Source=C:/Dental.mdb;")
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}

	if mongoDb == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, config.GetInstance().DbConnectionString)

		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}

		defer cancel()
		err = client.Ping(ctx, readpref.Primary())

		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}

		mongoDb = client
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
	coll := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.DiagnosisCollection)

	var items []interface{}

	for _, t := range *diagnosis {
		items = append(items, t)
	}

	_, err := coll.InsertMany(context.Background(), items)

	return err
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
	coll := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.ManipulationsCollection)

	var items []interface{}

	for _, t := range *manipulations {
		items = append(items, t)
	}

	_, err := coll.InsertMany(context.Background(), items)

	return err
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
	coll := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.ToothStatusesCollection)

	var items []interface{}

	for _, t := range *toothStatuses {
		items = append(items, t)
	}

	_, err := coll.InsertMany(context.Background(), items)

	return err
}

// pushDentist ...
func pushDentist() (primitive.ObjectID, error) {
	coll := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.DentistCollection)

	err := coll.Drop(context.Background())

	if err != nil {
		return primitive.NilObjectID, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("1"), 8)
	if err != nil {
		return primitive.NilObjectID, err
	}

	dentist := nosql.Dentist{
		Email:    "a@a.a",
		Name:     "Маргарита Търпанова",
		Password: hashedPassword,
	}

	newDentist, err := coll.InsertOne(context.Background(), dentist)
	if err != nil {
		return primitive.NilObjectID, err
	}

	newDentistID := newDentist.InsertedID.(primitive.ObjectID)

	return newDentistID, nil
}

// seedPatients ..
func seedPatients() (*[]*nosql.Patient, error) {

	patients := make([]*nosql.Patient, 0)
	m := make(map[int]*nosql.Patient)

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
			patient = &nosql.Patient{
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

			patient.Teeth = make([]*nosql.Tooth, 0)
			m[patientID] = patient
		}

		tooth, err := patient.GetToothByNo(toothNo)

		if err != nil {
			return nil, err
		}

		if tooth == nil {
			tooth = &nosql.Tooth{
				ToothNo:       strconv.Itoa(toothNo),
				Diagnosis:     make([]*nosql.ToothOperation, 0),
				Manipulations: make([]*nosql.ToothOperation, 0),
			}

			patient.Teeth = append(patient.Teeth, tooth)
		}

		var recordID uuid.UUID
		recordID, err = uuid.NewV4()

		if err != nil {
			return nil, err
		}

		operation := nosql.ToothOperation{
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
			patient = &nosql.Patient{
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

			patient.Teeth = make([]*nosql.Tooth, 0)
			m[patientID] = patient
		}

		tooth, err := patient.GetToothByNo(toothNo)

		if err != nil {
			return nil, err
		}

		if tooth == nil {
			tooth = &nosql.Tooth{
				ToothNo:       strconv.Itoa(toothNo),
				Diagnosis:     make([]*nosql.ToothOperation, 0),
				Manipulations: make([]*nosql.ToothOperation, 0),
			}

			patient.Teeth = append(patient.Teeth, tooth)
		}

		var recordID uuid.UUID
		recordID, err = uuid.NewV4()

		if err != nil {
			return nil, err
		}

		operation := nosql.ToothOperation{
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
func pushPatients(dentistID primitive.ObjectID, patients *[]*nosql.Patient) error {

	coll := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.PatientCollection)

	var items []interface{}

	for _, t := range *patients {

		t.Dentists = make([]primitive.ObjectID, 0)
		t.Dentists = append(t.Dentists, dentistID)

		items = append(items, *t)
	}

	r, err := coll.InsertMany(context.Background(), items)
	if err != nil {
		return err
	}

	ids := r.InsertedIDs

	dentistsCollection := mongoDb.Database(nosql.MongoDbSchema.DatabaseName).Collection(nosql.MongoDbSchema.DentistCollection)
	dentistFilter := bson.M{"_id": dentistID}

	_, err = dentistsCollection.UpdateOne(context.Background(), dentistFilter, bson.M{"$set": bson.M{"patients": ids}})

	return err
}
