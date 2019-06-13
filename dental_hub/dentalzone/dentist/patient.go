package dentist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	err = repo.UpdatePatientProfile(&newParient)

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

	input, err := json.Marshal(newParient)

	fmt.Println(string(input))
	if err != nil {
		return err
	}

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	newParient.RegistrationDate = time.Now()

	patientID, err := repo.CreatePatientProfile(&newParient, dentistID)

	if err != nil {
		return err
	}

	output, _ := json.Marshal(patientID)
	fmt.Fprintf(w, string(output))

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
	err = repo.RemovePatientProfile(patientID, dentistID)

	if err != nil {
		return err
	}

	return nil
}

// Download ...
func Download(w http.ResponseWriter, r *http.Request) error {

	reader, err := repo.GetImage("", "")
	if err != nil {
		return err
	}

	var nBytes int64
	nBytes, err = io.Copy(w, reader)

	fmt.Println("", nBytes)

	return err
}

// Upload ...
func Upload(w http.ResponseWriter, r *http.Request) error {

	patientID := r.URL.Query().Get("id")
	fmt.Println(patientID)

	file, header, err := r.FormFile("fileUpload")
	if err != nil {
		return err
	}

	tags := r.FormValue("tags")

	fmt.Println(header.Filename)
	fmt.Println(tags)

	/*
		defer file.Close()

		// copy example
		f, err := os.OpenFile("./"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		defer f.Close()

		io.Copy(f, file)
		return nil
	*/

	repo := repository.Repository

	tags2 := strings.Fields(tags)
	err = repo.InsertImage(patientID, file, tags2, header.Size)

	if err != nil {
		return err
	}

	return nil
}
