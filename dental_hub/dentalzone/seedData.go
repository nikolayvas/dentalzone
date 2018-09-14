package dentalzone

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "dental_hub/models"

	"dental_hub/repository"
)

var repo = repository.Repository

// SeedDataViewModel model
type SeedDataViewModel struct {
	Diagnosises   []m.Diagnosis    `json:"diagnosisList"`
	Manipulations []m.Manipulation `json:"manipulationList"`
	ToothStatuses []m.ToothStatus  `json:"toothStatusList"`
}

// SeedData endpoint
func SeedData(w http.ResponseWriter, r *http.Request) error {

	diagnosis, err := repo.SeedDiagnosis()
	if err != nil {
		return err
	}

	manipulations, err := repo.SeedManipulations()
	if err != nil {
		return err
	}

	toothStatuses, err := repo.SeedToothStatuses()
	if err != nil {
		return err
	}

	data := SeedDataViewModel{Diagnosises: diagnosis, Manipulations: manipulations, ToothStatuses: toothStatuses}
	output, _ := json.Marshal(data)
	fmt.Fprintf(w, string(output))

	return nil
}
