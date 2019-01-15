package dentist

import (
	"dental_hub/core"
	m "dental_hub/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SeedAppointments endpoint
func SeedAppointments(w http.ResponseWriter, r *http.Request) error {

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	date := time.Now()

	appointments, err := repo.GetAppointments(dentistID, date)

	if err != nil {
		return err
	}

	output, _ := json.Marshal(appointments)
	fmt.Fprintf(w, string(output))

	return nil
}

// UpdateAppointments endpoint
func UpdateAppointments(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var appointments m.Appointments
	err := decoder.Decode(&appointments)
	if err != nil {
		return err
	}

	output, err := json.Marshal(appointments)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	err = repo.UpdateAppointments(dentistID, appointments.Date, &appointments.Appointments)

	if err != nil {
		return err
	}

	return nil
}
