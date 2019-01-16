package nosql

import (
	"context"
	ex "dental_hub/exceptions"
	m "dental_hub/models"
	u "dental_hub/utils"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/satori/go.uuid"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func contextWithTimeout(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

// Dentist read model
type Dentist struct {
	ID       primitive.ObjectID     `bson:"_id,omitempty"`
	Email    string                 `bson:"email"`
	Name     string                 `bson:"name"`
	Password []byte                 `bson:"password"`
	Patients [](primitive.ObjectID) `bson:"patients"`
}

// Tooth status
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

// Patient read/write model
type Patient struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	FirstName        string               `bson:"firstName"`
	MiddleName       string               `bson:"middleName"`
	LastName         string               `bson:"lastName"`
	Email            string               `bson:"email"`
	Address          string               `bson:"address"`
	PhoneNumber      string               `bson:"phoneNumber"`
	GeneralInfo      string               `bson:"generalInfo"`
	RegistrationDate time.Time            `bson:"registrationDate"`
	Dentists         []primitive.ObjectID `bson:"dentists"`
	Teeth            []*Tooth             `bson:"teeth"`
}

// SignUp dentist write model
type SignUp struct {
	Email          string    `bson:"email"`
	UserName       string    `bson:"userName"`
	Password       []byte    `bson:"password"`
	VerificationID string    `bson:"verificationID"`
	ExpirationDate time.Time `bson:"expirationDate"`
}

// ResetPassword write model
type ResetPassword struct {
	DentistID      primitive.ObjectID `bson:"dentistId"`
	Code           string             `bson:"code"`
	ExpirationDate time.Time          `bson:"expirationDate"`
}

// RegisterDentist registers new user and returns verification code
func (r Repository) RegisterDentist(email string, userName string, password []byte) (string, error) {
	dentist := Dentist{}

	dentistsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)
	filter := bson.M{"email": email}
	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	// check if already have registered user with such email
	err := dentistsCollection.FindOne(ctx, filter).Decode(&dentist)

	switch {
	case err == mongo.ErrNoDocuments:
		// Do nothing, the email is available for next registration
	case err != nil:
		return "", err
	default:
		return "", ex.ErrAlreadyExists
	}

	// if already have unprocessed requests for registration with such email, remove them
	signUpCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistSignUpCollection)
	_, err = signUpCollection.DeleteMany(ctx, filter)

	if err != nil {
		return "", err
	}

	var verificationID uuid.UUID
	verificationID, err = uuid.NewV4()

	if err != nil {
		return "", err
	}

	verificationIDToString := verificationID.String()

	doc := SignUp{
		UserName:       userName,
		Email:          email,
		Password:       password,
		VerificationID: verificationIDToString,
		ExpirationDate: time.Now().Add(3 * time.Hour),
	}

	// Store user information in temporary table and wait for email verification
	_, err = signUpCollection.InsertOne(ctx, doc)

	if err != nil {
		return "", err
	}

	return verificationIDToString, nil
}

// ActivateDentist activates alredy registered user
func (r Repository) ActivateDentist(id string) error {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	signUp := SignUp{}

	signUpCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistSignUpCollection)
	filter := bson.M{"verificationID": id, "expirationDate": bson.M{"$gte": time.Now()}}

	// check if exists such none expired code
	err := signUpCollection.FindOne(ctx, filter).Decode(&signUp)

	if err != nil {
		switch {
		case err == mongo.ErrNoDocuments:
			return ex.ErrNotSuch
		default:
			return err
		}
	}

	doc := Dentist{
		Name:     signUp.UserName,
		Email:    signUp.Email,
		Password: signUp.Password,
	}

	dentistCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)
	_, err = dentistCollection.InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

// Login returns user details
func (r Repository) Login(email string) (*m.Login, error) {

	_, _ = r.SeedToothStatuses()

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	dentist := Dentist{}

	dentistCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)
	filter := bson.M{"email": email}

	err := dentistCollection.FindOne(ctx, filter).Decode(&dentist)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ex.ErrNotSuch
	case err != nil:
		return nil, err
	default:
		login := m.Login{
			ID:       dentist.ID.Hex(),
			Email:    dentist.Email,
			Password: dentist.Password,
			Name:     dentist.Name,
		}
		return &login, nil
	}
}

// AddPasswordResetConfirmationCode insertс new confirmation code in db
func (r Repository) AddPasswordResetConfirmationCode(email string, code string) error {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	dentistResetPasswordCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistResetPasswordCollection)
	filter := bson.M{"email": email}

	_, err := dentistResetPasswordCollection.DeleteMany(ctx, filter)

	if err != nil {
		return err
	}

	dentist := Dentist{}
	dentistCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)

	err = dentistCollection.FindOne(ctx, filter).Decode(&dentist)

	if err != nil {
		return err
	}

	resetPassword := ResetPassword{
		DentistID:      dentist.ID,
		Code:           code,
		ExpirationDate: time.Now().Add(3 * time.Hour),
	}

	_, err = dentistResetPasswordCollection.InsertOne(ctx, resetPassword)
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword resets user password
func (r Repository) ResetPassword(hashedPassword []byte, email string, code string) error {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	dentist := Dentist{}

	dentistCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)
	dentistFilter := bson.M{"email": email}

	err := dentistCollection.FindOne(ctx, dentistFilter).Decode(&dentist)

	if err != nil {
		return err
	}

	dentistResetPasswordCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistResetPasswordCollection)
	passwordResetFilter := bson.M{"dentistId": dentist.ID, "code": code, "expirationDate": bson.M{"$gte": time.Now()}}

	err = dentistResetPasswordCollection.FindOne(ctx, passwordResetFilter).Decode(&dentist)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ex.ErrNotSuch
		}

		return err
	}

	dentist.Password = hashedPassword
	_, err = dentistCollection.UpdateOne(ctx, dentistFilter, bson.M{"$set": bson.M{"password": hashedPassword}})

	if err != nil {
		return err
	}

	return nil
}

// SeedDiagnosis seeds diagnosis
func (r Repository) SeedDiagnosis() (*[]m.Diagnosis, error) {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	coll := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DiagnosisCollection)

	cursor, err := coll.Find(
		ctx,
		bson.D{},
	)

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	diagnosisList := make([]m.Diagnosis, 0)

	for cursor.Next(ctx) {
		var diagnosis m.Diagnosis

		err = cursor.Decode(&diagnosis)

		if err != nil {
			return nil, err
		}

		diagnosisList = append(diagnosisList, diagnosis)
	}

	/*
		_, err := coll.InsertMany(
			ctx,
			[]interface{}{
				bson.D{
					{"id", bsonx.Int32(1)},
					{"diagnosisName", bsonx.String("Caries superficialis")},
					{"changeStatus", bsonx.Int32(3)},
				},
				bson.D{
					{"id", bsonx.Int32(2)},
					{"diagnosisName", bsonx.String("Caries media")},
					{"changeStatus", bsonx.Int32(3)},
				},
				bson.D{
					{"id", bsonx.Int32(3)},
					{"diagnosisName", bsonx.String("Caries profunda")},
					{"changeStatus", bsonx.Int32(3)},
				},
				bson.D{
					{"id", bsonx.Int32(4)},
					{"diagnosisName", bsonx.String("Pulpitis acuta serosa partialis")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(5)},
					{"diagnosisName", bsonx.String("Pulpitis acuta serosa totalis")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(6)},
					{"diagnosisName", bsonx.String("Pulpitis acuta purulenta partialis")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(7)},
					{"diagnosisName", bsonx.String("Pulpitis acuta purulenta totalis")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(8)},
					{"diagnosisName", bsonx.String("Pulpitis acuta gangraenosa")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(9)},
					{"diagnosisName", bsonx.String("Pulpitis chronica ulcerosa")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(10)},
					{"diagnosisName", bsonx.String("Pulpitis chronica granulomatosa")},
					{"changeStatus", bsonx.Int32(4)},
				},
				bson.D{
					{"id", bsonx.Int32(11)},
					{"diagnosisName", bsonx.String("Periodontitis acuta serosa")},
					{"changeStatus", bsonx.Int32(5)},
				},
				bson.D{
					{"id", bsonx.Int32(12)},
					{"diagnosisName", bsonx.String("Periodontitis acuta purulenta")},
					{"changeStatus", bsonx.Int32(5)},
				},
				bson.D{
					{"id", bsonx.Int32(13)},
					{"diagnosisName", bsonx.String("Periodontitis chronica fibrosa")},
					{"changeStatus", bsonx.Int32(5)},
				},
				bson.D{
					{"id", bsonx.Int32(14)},
					{"diagnosisName", bsonx.String("Periodontitis chronica granulomatosa progressiva (cum-sine fistula)")},
					{"changeStatus", bsonx.Int32(5)},
				},
				bson.D{
					{"id", bsonx.Int32(15)},
					{"diagnosisName", bsonx.String("Periodontitis chronica granulomatosa localisata")},
					{"changeStatus", bsonx.Int32(5)},
				},
				bson.D{
					{"id", bsonx.Int32(16)},
					{"diagnosisName", bsonx.String("Radix")},
					{"changeStatus", bsonx.Int32(8)},
				},
				bson.D{
					{"id", bsonx.Int32(17)},
					{"diagnosisName", bsonx.String("Extractio")},
					{"changeStatus", bsonx.Int32(2)},
				},
				bson.D{
					{"id", bsonx.Int32(18)},
					{"diagnosisName", bsonx.String("Building")},
					{"changeStatus", bsonx.Int32(7)},
				},
				bson.D{
					{"id", bsonx.Int32(19)},
					{"diagnosisName", bsonx.String("Obturatio")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(20)},
					{"diagnosisName", bsonx.String("Corona")},
					{"changeStatus", bsonx.Int32(9)},
				},
				bson.D{
					{"id", bsonx.Int32(21)},
					{"diagnosisName", bsonx.String("Dens intact")},
					{"changeStatus", bsonx.Int32(1)},
				},
			})
	*/

	return &diagnosisList, nil
}

// SeedManipulations seeds manipulations
func (r Repository) SeedManipulations() (*[]m.Manipulation, error) {
	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	coll := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.ManipulationsCollection)

	/*
		_, err := coll.InsertMany(
			ctx,
			[]interface{}{
				bson.D{
					{"id", bsonx.Int32(1)},
					{"manipulationName", bsonx.String("Caries himiopolimer")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(2)},
					{"manipulationName", bsonx.String("Caries fotopolimer")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(3)},
					{"manipulationName", bsonx.String("Pulpit one canal")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(4)},
					{"manipulationName", bsonx.String("Pulpit two canalis")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(5)},
					{"manipulationName", bsonx.String("Pulpit three canalis")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(6)},
					{"manipulationName", bsonx.String("Periodontit one canal")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(7)},
					{"manipulationName", bsonx.String("Periodontit two canalis")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(8)},
					{"manipulationName", bsonx.String("Periodontit two canalis")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(9)},
					{"manipulationName", bsonx.String("Periodontit three canalis")},
					{"changeStatus", bsonx.Int32(6)},
				},
				bson.D{
					{"id", bsonx.Int32(10)},
					{"manipulationName", bsonx.String("Extractio one radix")},
					{"changeStatus", bsonx.Int32(2)},
				},
				bson.D{
					{"id", bsonx.Int32(11)},
					{"manipulationName", bsonx.String("Extractio three radix")},
					{"changeStatus", bsonx.Int32(2)},
				},
				bson.D{
					{"id", bsonx.Int32(12)},
					{"manipulationName", bsonx.String("Extractio madrec")},
					{"changeStatus", bsonx.Int32(2)},
				},
				bson.D{
					{"id", bsonx.Int32(13)},
					{"manipulationName", bsonx.String("Anestsia")},
					{"changeStatus", bsonx.Int32(1)},
				},
				bson.D{
					{"id", bsonx.Int32(14)},
					{"manipulationName", bsonx.String("Building")},
					{"changeStatus", bsonx.Int32(7)},
				},
				bson.D{
					{"id", bsonx.Int32(15)},
					{"manipulationName", bsonx.String("Cimentirane na stara corona")},
					{"changeStatus", bsonx.Int32(9)},
				},
				bson.D{
					{"id", bsonx.Int32(16)},
					{"manipulationName", bsonx.String("Svaliane na stara corona")},
					{"changeStatus", bsonx.Int32(1)},
				},
				bson.D{
					{"id", bsonx.Int32(17)},
					{"manipulationName", bsonx.String("Korona 'Adapta'")},
					{"changeStatus", bsonx.Int32(9)},
				},
				bson.D{
					{"id", bsonx.Int32(18)},
					{"manipulationName", bsonx.String("Korona 'Metalokeramika'")},
					{"changeStatus", bsonx.Int32(9)},
				},
			})
	*/

	cursor, err := coll.Find(
		ctx,
		bson.D{},
	)

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	manipulationList := make([]m.Manipulation, 0)

	for cursor.Next(ctx) {
		var manipilation m.Manipulation

		err = cursor.Decode(&manipilation)

		if err != nil {
			return nil, err
		}

		manipulationList = append(manipulationList, manipilation)
	}

	return &manipulationList, nil
}

// SeedToothStatuses seeds tooth statuses
func (r Repository) SeedToothStatuses() (*[]m.ToothStatus, error) {
	ctx := context.Background()

	coll := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.ToothStatusesCollection)

	/*
		_, err := coll.InsertMany(
			ctx,
			[]interface{}{
				bson.D{
					{"id", bsonx.Int32(1)},
					{"status", bsonx.String("Unchenged")},
				},
				bson.D{
					{"id", bsonx.Int32(2)},
					{"status", bsonx.String("Extractio")},
				},
				bson.D{
					{"id", bsonx.Int32(3)},
					{"status", bsonx.String("Caries")},
				},
				bson.D{
					{"id", bsonx.Int32(4)},
					{"status", bsonx.String("Pulpitis")},
				},
				bson.D{
					{"id", bsonx.Int32(5)},
					{"status", bsonx.String("Periodontitis")},
				},
				bson.D{
					{"id", bsonx.Int32(6)},
					{"status", bsonx.String("Obturacio")},
				},
				bson.D{
					{"id", bsonx.Int32(7)},
					{"status", bsonx.String("Building")},
				},
				bson.D{
					{"id", bsonx.Int32(8)},
					{"status", bsonx.String("Radix")},
				},
				bson.D{
					{"id", bsonx.Int32(9)},
					{"status", bsonx.String("Corona")},
				},
			})
	*/
	cursor, err := coll.Find(
		ctx,
		bson.D{},
	)

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	statusesList := make([]m.ToothStatus, 0)

	for cursor.Next(ctx) {
		var status m.ToothStatus

		err = cursor.Decode(&status)

		if err != nil {
			return nil, err
		}

		statusesList = append(statusesList, status)
	}

	return &statusesList, nil
}

// CreatePatientProfile updates patient
func (r Repository) CreatePatientProfile(newParient m.Patient, dentistID string) error {

	dentist := Dentist{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	dentistsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)

	hex, err := primitive.ObjectIDFromHex(dentistID)
	dentistFilter := bson.M{"_id": hex}

	err = dentistsCollection.FindOne(ctx, dentistFilter).Decode(&dentist)

	if err != nil {
		return err
	}

	patientDoc := Patient{
		Address:          newParient.Address,
		Dentists:         []primitive.ObjectID{hex},
		Email:            newParient.Email,
		FirstName:        newParient.FirstName,
		MiddleName:       newParient.MiddleName,
		GeneralInfo:      newParient.GeneralInfo,
		LastName:         newParient.LastName,
		PhoneNumber:      newParient.PhoneNumber,
		RegistrationDate: newParient.RegistrationDate,
	}

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)

	patient, err := patientCollection.InsertOne(ctx, patientDoc)

	if err != nil {
		return err
	}

	newPatientID := patient.InsertedID.(primitive.ObjectID)

	newPatientsByDentist := append(dentist.Patients, newPatientID)

	_, err = dentistsCollection.UpdateOne(ctx, dentistFilter, bson.M{"$set": bson.M{"patients": newPatientsByDentist}})

	if err != nil {
		return err
	}

	return nil
}

// UpdatePatientProfile updates patient
func (r Repository) UpdatePatientProfile(patient m.Patient) error {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(patient.ID)

	if err != nil {
		return err
	}

	_, err = patientCollection.UpdateOne(ctx, bson.M{"_id": hex}, bson.M{"$set": bson.M{
		"address":     patient.Address,
		"firstName":   patient.FirstName,
		"middleName":  patient.MiddleName,
		"lastName":    patient.LastName,
		"phoneNumber": patient.PhoneNumber,
		"generalInfo": patient.GeneralInfo,
	}})

	return err
}

// GetPatients returns patients
func (r Repository) GetPatients(dentistID string) (*[]m.Patient, error) {

	dentist := Dentist{}
	ctx := context.Background()

	dentistsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)

	hex, err := primitive.ObjectIDFromHex(dentistID)

	if err != nil {
		return nil, err
	}

	dentistFilter := bson.M{"_id": hex}

	err = dentistsCollection.FindOne(ctx, dentistFilter).Decode(&dentist)

	if err != nil {
		return nil, err
	}

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	patientsFilter := bson.M{"_id": bson.M{"$in": dentist.Patients}}

	cursor, err := patientCollection.Find(ctx, patientsFilter)

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	patients := make([]m.Patient, 0)

	for cursor.Next(ctx) {
		var patient Patient

		err = cursor.Decode(&patient)

		if err != nil {
			return nil, err
		}

		mongoID := patient.ID.Hex()
		patients = append(patients, m.Patient{
			ID:               mongoID,
			Address:          patient.Address,
			Email:            patient.Email,
			FirstName:        patient.FirstName,
			MiddleName:       patient.MiddleName,
			LastName:         patient.LastName,
			GeneralInfo:      patient.GeneralInfo,
			PhoneNumber:      patient.PhoneNumber,
			RegistrationDate: patient.RegistrationDate,
		})
	}

	return &patients, nil
}

// RemovePatientProfile removes the patient from the list of patients for the dentist
func (r Repository) RemovePatientProfile(patientID string, dentistID string) error {
	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(patientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return err
	}

	//TODO If the patient has multiple dentists, remove just the current one
	_, err = patientCollection.UpdateOne(ctx, patientFilter, bson.M{"$set": bson.M{"dentists": nil}})

	dentistsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)

	var dentist Dentist
	var newPatientsByDentist [](primitive.ObjectID)

	for _, dentistID := range patient.Dentists {
		dentist = Dentist{}
		err = dentistsCollection.FindOne(ctx, bson.M{"_id": dentistID}).Decode(&dentist)

		if err != nil {
			return err
		}

		for i, id := range dentist.Patients {
			if id.Hex() == patientID {
				newPatientsByDentist = append(dentist.Patients[:i], dentist.Patients[i+1:]...)
				break
			}
		}

		_, err = dentistsCollection.UpdateOne(ctx, bson.M{"_id": dentistID}, bson.M{"$set": bson.M{"patients": newPatientsByDentist}})

		if err != nil {
			return err
		}
	}

	return nil
}

// GetTeethData returns teeth data per patient
func (r Repository) GetTeethData(patientID string) (*m.TeethData, error) {

	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(patientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return nil, err
	}

	var result = m.TeethData{
		Diagnosis:     []m.ToothAction{},
		Manipulations: []m.ToothAction{},
	}

	for _, tooth := range patient.Teeth {
		for _, diagnosis := range tooth.Diagnosis {
			result.Diagnosis = append(result.Diagnosis, m.ToothAction{
				ID:        diagnosis.RecordID,
				ActionID:  diagnosis.OperationID,
				PatientID: patientID,
				ToothNo:   tooth.ToothNo,
				Date:      diagnosis.Date,
			})
		}

		for _, manipulation := range tooth.Manipulations {
			result.Manipulations = append(result.Manipulations, m.ToothAction{
				ID:        manipulation.RecordID,
				ActionID:  manipulation.OperationID,
				PatientID: patientID,
				ToothNo:   tooth.ToothNo,
				Date:      manipulation.Date,
			})
		}
	}

	return &result, nil
}

// AddToothManipulation adds manipulation
func (r Repository) AddToothManipulation(manipulation m.ToothAction) error {
	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(manipulation.PatientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return err
	}

	_, tooth := FindTooth(patient.Teeth, manipulation.ToothNo)

	if tooth != nil {
		var newList = append(tooth.Manipulations, &ToothOperation{
			Date:        time.Now(),
			RecordID:    manipulation.ID,
			OperationID: manipulation.ActionID,
		})

		tooth.Manipulations = newList
	} else {
		tooth = &Tooth{
			ToothNo: manipulation.ToothNo,
			Manipulations: []*ToothOperation{
				&ToothOperation{
					Date:        time.Now(),
					RecordID:    manipulation.ID,
					OperationID: manipulation.ActionID,
				},
			},
		}

		patient.Teeth = append(patient.Teeth, tooth)
	}

	_, err = patientCollection.UpdateOne(ctx, patientFilter, bson.M{"$set": bson.M{"teeth": patient.Teeth}})

	return err
}

// RemoveToothManipulation removes manipulation
func (r Repository) RemoveToothManipulation(manipulation m.ToothAction) error {
	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(manipulation.PatientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return err
	}

	_, tooth := FindTooth(patient.Teeth, manipulation.ToothNo)

	if tooth != nil {
		i, operation := FindToothOperation(tooth.Manipulations, manipulation.ID)

		if operation != nil {
			tooth.Manipulations[i] = &ToothOperation{
				RecordID:    operation.RecordID,
				OperationID: operation.OperationID,
				Date:        operation.Date,
				IsDeleted:   u.RefBool(true),
				DateDeleted: u.RefTime(time.Now()),
			}
		} else {
			//TODO ??
		}

	} else {
		// TODO ???
	}

	_, err = patientCollection.UpdateOne(ctx, patientFilter, bson.M{"$set": bson.M{"teeth": patient.Teeth}})

	return err
}

// AddToothDiagnosis adds diagnosis
func (r Repository) AddToothDiagnosis(diagnosis m.ToothAction) error {
	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(diagnosis.PatientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return err
	}

	_, tooth := FindTooth(patient.Teeth, diagnosis.ToothNo)

	if tooth != nil {
		var newList = append(tooth.Diagnosis, &ToothOperation{
			Date:        time.Now(),
			RecordID:    diagnosis.ID,
			OperationID: diagnosis.ActionID,
		})

		tooth.Diagnosis = newList
	} else {
		tooth = &Tooth{
			ToothNo: diagnosis.ToothNo,
			Diagnosis: []*ToothOperation{
				&ToothOperation{
					Date:        time.Now(),
					RecordID:    diagnosis.ID,
					OperationID: diagnosis.ActionID,
				},
			},
		}

		patient.Teeth = append(patient.Teeth, tooth)
	}

	_, err = patientCollection.UpdateOne(ctx, patientFilter, bson.M{"$set": bson.M{"teeth": patient.Teeth}})

	return err
}

// RemoveToothDiagnosis removes diagnosis
func (r Repository) RemoveToothDiagnosis(diagnosis m.ToothAction) error {
	patient := Patient{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	patientCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.PatientCollection)
	hex, err := primitive.ObjectIDFromHex(diagnosis.PatientID)
	patientFilter := bson.M{"_id": hex}

	err = patientCollection.FindOne(ctx, patientFilter).Decode(&patient)

	if err != nil {
		return err
	}

	_, tooth := FindTooth(patient.Teeth, diagnosis.ToothNo)

	if tooth != nil {
		i, operation := FindToothOperation(tooth.Diagnosis, diagnosis.ID)

		if operation != nil {
			tooth.Diagnosis[i] = &ToothOperation{
				RecordID:    operation.RecordID,
				OperationID: operation.OperationID,
				Date:        operation.Date,
				IsDeleted:   u.RefBool(true),
				DateDeleted: u.RefTime(time.Now()),
			}
		} else {
			//TODO ??
		}

	} else {
		// TODO ???
	}

	_, err = patientCollection.UpdateOne(ctx, patientFilter, bson.M{"$set": bson.M{"teeth": patient.Teeth}})

	return err
}

// InvitePatient resets patient password
func (r Repository) InvitePatient(dentistID string, patientEmail string) (string, error) {

	return "", nil

	//TODO what if already assigned
}

// ActivateInvitation activates invitation
func (r Repository) ActivateInvitation(activationID string) error {
	return nil
}

// GetDentist returns dentist data
func (r Repository) GetDentist(dentistID string) (*m.Dentist, error) {

	dentist := Dentist{}

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	dentistCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.DentistCollection)
	hex, err := primitive.ObjectIDFromHex(dentistID)

	if err != nil {
		return nil, err
	}

	err = dentistCollection.FindOne(ctx, bson.M{"_id": hex}).Decode(&dentist)

	if err != nil {
		return nil, err
	}

	return &m.Dentist{Name: dentist.Name, Email: dentist.Email}, err
}

// FindTooth in the collection by toothNo
func FindTooth(teeth []*Tooth, toothNo string) (int, *Tooth) {
	if teeth == nil {
		return -1, nil
	}

	for i, t := range teeth {
		if t.ToothNo == toothNo {
			return i, t
		}
	}

	return -1, nil
}

// FindToothOperation in the collection by RecordNo
func FindToothOperation(operations []*ToothOperation, recordID string) (int, *ToothOperation) {
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

// GetAppointments returns appointments for day
func (r Repository) GetAppointments(dentistID string, date time.Time) (*[]m.Appointment, error) {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	scheduleCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.ScheduleCollection)
	scheduleFilter := bson.M{"dentistID": dentistID, "day": rounded}

	appointments := m.Appointments{}
	err := scheduleCollection.FindOne(ctx, scheduleFilter).Decode(&appointments)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &appointments.Appointments, err
	}
}

type result interface{}

// UpdateAppointments updates appointments for day
func (r Repository) UpdateAppointments(dentistID string, date time.Time, appointments *[]m.Appointment) error {

	ctx, cancel := contextWithTimeout(10)
	defer cancel()

	rounded := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	scheduleCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.ScheduleCollection)
	scheduleFilter := bson.M{"dentistID": dentistID, "day": rounded}

	t := true
	opt := options.UpdateOptions{Upsert: &t}
	_, err := scheduleCollection.UpdateOne(ctx, scheduleFilter, bson.M{"$set": bson.M{"appointments": *appointments}}, &opt)

	return err
}
