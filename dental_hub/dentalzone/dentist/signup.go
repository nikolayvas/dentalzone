package dentist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	config "dental_hub/configuration"
	ex "dental_hub/exceptions"
	i "dental_hub/infrastructure"
)

// SignUpViewModel model
type SignUpViewModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// SignUpRegister endpoint
func SignUpRegister(w http.ResponseWriter, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)

	var dentist SignUpViewModel
	err := decoder.Decode(&dentist)
	if err != nil {
		return err
	}

	if len(dentist.Password) == 0 {
		http.Error(w, "Password is incorrect!", http.StatusNotAcceptable)
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dentist.Password), 8)
	if err != nil {
		return err
	}

	verificationID, err := repo.RegisterDentist(dentist.Email, dentist.Name, hashedPassword)

	if err != nil {
		switch {
		case err == ex.ErrAlreadyExists:
			http.Error(w, "Already have reqistered user with such email!", http.StatusNotAcceptable)
			return nil
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
	}

	hyperlink := fmt.Sprintf("%s?id=%s", config.GetInstance().DentistProfileActivateURI, *verificationID)
	body := "Your registration is almost ready. Please click over the link to activate your profile: \r\n" +
		fmt.Sprintf("<a href=\"%s\">%s</a>", hyperlink, hyperlink)

	mail := &i.Mail{
		To:      []string{dentist.Email},
		Subject: "DentalZone Registration",
		Body:    body}

	err = i.SendMail(mail)

	if err != nil {
		return err
	}

	return nil
}

// SignUpActivate endpoint
func SignUpActivate(w http.ResponseWriter, r *http.Request) error {

	id := r.URL.Query().Get("id")

	err := repo.ActivateDentist(id)

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
