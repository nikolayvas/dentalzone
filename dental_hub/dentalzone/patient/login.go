package patient

import (
	"encoding/json"
	"fmt"
	"net/http"

	ex "dental_hub/exceptions"

	"golang.org/x/crypto/bcrypt"

	"dental_hub/core"
	u "dental_hub/utils"
)

// ForgotPasswordViewModel model
type ForgotPasswordViewModel struct {
	Email string `json:"email"`
}

// ResetPasswordViewModel model
type ResetPasswordViewModel struct {
	Email    string `json:"email"`
	Password string `json:"newPassword"`
	Code     string `json:"confirmationCode"`
}

// LoginViewModel model
type LoginViewModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginDto dto
type LoginDto struct {
	Email string  `json:"email"`
	Name  string  `json:"name"`
	Token *string `json:"token"`
}

// Login endpoint
func Login(w http.ResponseWriter, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)

	var user LoginViewModel
	err := decoder.Decode(&user)
	if err != nil {
		return err
	}

	login, err := repo.LoginPatient(user.Email)

	switch {
	case err == ex.ErrNotSuch:
		http.Error(w, "No such user!", http.StatusNotFound)
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		if err = bcrypt.CompareHashAndPassword(login.Password, []byte(user.Password)); err != nil {
			w.WriteHeader(http.StatusNotFound)

			return nil
		}

		token, err := core.GenerateJwt(login.ID)
		if err != nil {
			return err
		}

		output, _ := json.Marshal(LoginDto{Email: login.Email, Name: login.Name, Token: token})
		fmt.Fprintf(w, string(output))
	}

	return nil
}

// SendPasswordResetConfirmationCode endpoint
func SendPasswordResetConfirmationCode(w http.ResponseWriter, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)

	var model ForgotPasswordViewModel
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}

	code := u.EncodeToString(6)

	err = repo.AddPatientPasswordResetConfirmationCode(model.Email, code)

	if err != nil {
		return err
	}

	//TODO Send email with the code
	return nil
}

// ResetPassword endpoint
func ResetPassword(w http.ResponseWriter, r *http.Request) error {

	decoder := json.NewDecoder(r.Body)

	var model ResetPasswordViewModel
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(model.Password), 8)
	if err != nil {
		return err
	}

	err = repo.ResetPatientPassword(hashedPassword, model.Email, model.Code)

	switch {
	case err == ex.ErrNotSuch:
		http.Error(w, "Unsuccessful reset! Possible reasons could be non existing email or expired or not matching verification code.", http.StatusNotAcceptable)
	case err != nil:
		return err
	default:
		return nil
	}

	return nil
}
