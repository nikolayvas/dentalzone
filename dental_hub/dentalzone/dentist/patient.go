package dentist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dental_hub/core"
	m "dental_hub/models"
	"dental_hub/repository"
)

// GetPatients endpoint
func GetPatients(w http.ResponseWriter, r *http.Request) error {

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	patients, err := repo.GetPatients(dentistID)

	if err != nil {
		return err
	}

	output, _ := json.Marshal(patients)
	fmt.Fprintf(w, string(output))

	return nil
}

// UpdatePatientProfile endpoint
func UpdatePatientProfile(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var newParient m.Patient
	err := decoder.Decode(&newParient)
	if err != nil {
		return err
	}

	output, err := json.Marshal(newParient)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	err = repo.UpdatePatientProfile(newParient)

	if err != nil {
		return err
	}

	return nil
}

// CreatePatientProfile endpoint
func CreatePatientProfile(w http.ResponseWriter, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)

	var newParient m.Patient
	err := decoder.Decode(&newParient)
	if err != nil {
		return err
	}

	output, err := json.Marshal(newParient)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	newParient.RegistrationDate = time.Now()
	err = repo.CreatePatientProfile(newParient, dentistID)

	if err != nil {
		return err
	}

	return nil
}

//RemovePatientProfile endpoint
func RemovePatientProfile(w http.ResponseWriter, r *http.Request) error {

	patientID := r.URL.Query().Get("id")

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	repo := repository.Repository
	err = repo.RemovePatientProfile(&patientID, dentistID)

	if err != nil {
		return err
	}

	return nil
}
