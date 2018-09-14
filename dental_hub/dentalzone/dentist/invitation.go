package dentist

import (
	"encoding/json"
	"fmt"
	"net/http"

	config "dental_hub/configuration"
	"dental_hub/core"
	ex "dental_hub/exceptions"
	i "dental_hub/infrastructure"
)

// EmailViewModel model
type EmailViewModel struct {
	Email string `json:"email"`
}

// InvitePatient endpoint
func InvitePatient(w http.ResponseWriter, r *http.Request) error {

	dentistID, err := core.ExtractJwtClaim(r, "dentistId")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)

	var model EmailViewModel
	err = decoder.Decode(&model)
	if err != nil {
		return err
	}

	invitationID, err := repo.InvitePatient(*dentistID, model.Email)

	if err != nil {
		switch {
		case err == ex.ErrNotSuch:
			http.Error(w, "No such patient registered", http.StatusNotAcceptable)
			return nil
		default:
			return err
		}
	}

	hyperlink := fmt.Sprintf("%s?id=%s", config.GetInstance().InvitationActivateURI, *invitationID)
	body := "Your have been invited from ******. Please click over the link to accept invitation if you want: \r\n" +
		fmt.Sprintf("<a href=\"%s\">%s</a>", hyperlink, hyperlink)

	mail := &i.Mail{
		To:      []string{model.Email},
		Subject: "DentalZone Invitation",
		Body:    body}

	err = i.SendMail(mail)

	if err != nil {
		return err
	}

	return nil
}

// InvitationActivate endpoint
func InvitationActivate(w http.ResponseWriter, r *http.Request) error {

	id := r.URL.Query().Get("id")

	err := repo.ActivateInvitation(id)

	if err != nil {
		switch {
		case err == ex.ErrNotSuch:
			http.Error(w, "There is not such activation code or expiration time has been reached!", http.StatusNotAcceptable)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}
