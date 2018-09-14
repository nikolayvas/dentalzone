package dentist

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "dental_hub/models"
)

// GetTeethData endpoint
func GetTeethData(w http.ResponseWriter, r *http.Request) error {

	patientID := r.URL.Query().Get("patientId")

	result, err := repo.GetTeethData(patientID)

	if err != nil {
		return err
	}

	output, _ := json.Marshal(result)
	fmt.Fprintf(w, string(output))

	return nil
}

// AddToothManipulation endpoint
func AddToothManipulation(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var newThoothManupulation m.ToothAction
	err := decoder.Decode(&newThoothManupulation)
	if err != nil {
		return err
	}

	output, err := json.Marshal(newThoothManupulation)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	err = repo.AddToothManipulation(newThoothManupulation)

	if err != nil {
		return err
	}

	return nil
}

// RemoveToothManipulation endpoint
func RemoveToothManipulation(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")

	err := repo.RemoveToothManipulation(idParam)

	if err != nil {
		return err
	}

	return nil
}

// AddToothDiagnosis endpoint
func AddToothDiagnosis(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var newThoothDiagnosis m.ToothAction
	err := decoder.Decode(&newThoothDiagnosis)
	if err != nil {
		return err
	}

	output, err := json.Marshal(newThoothDiagnosis)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	err = repo.AddToothDiagnosis(newThoothDiagnosis)

	if err != nil {
		return err
	}

	return nil
}

// RemoveToothDiagnosis endpoint
func RemoveToothDiagnosis(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")

	err := repo.RemoveToothDiagnosis(idParam)

	if err != nil {
		return err
	}

	return nil
}
