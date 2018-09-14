package dentist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	c "dental_hub/core"
	m "dental_hub/models"
)

// Get return dentist logged in details
func Get(w http.ResponseWriter, r *http.Request) error {

	dentistID, err := c.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		//user is not logged in
		output, _ := json.Marshal(m.Dentist{})
		fmt.Fprintf(w, string(output))

		return nil
	}

	dentist, err := repo.GetDentist(dentistID)

	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "No such user!", http.StatusNotFound)
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		output, _ := json.Marshal(dentist)
		fmt.Fprintf(w, string(output))
	}

	return nil
}
