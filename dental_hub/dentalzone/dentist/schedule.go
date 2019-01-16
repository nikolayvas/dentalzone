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

	date := r.URL.Query().Get("day")
	time, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return err
	}

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	appointments, err := repo.GetAppointments(dentistID, time)

	if err != nil {
		return err
	}

	output, _ := json.Marshal(&appointments)
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

	err = repo.UpdateAppointments(dentistID, appointments.Day, &appointments.Appointments)

	if err != nil {
		return err
	}

	return nil
}
