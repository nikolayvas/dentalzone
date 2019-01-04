package dentalzone

import (
	"net/http"

	"github.com/gorilla/mux"

	c "dental_hub/core"

	d "dental_hub/dentalzone/dentist"
	p "dental_hub/dentalzone/patient"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc c.AppHandler
	Securitize  bool
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = c.Logger(handler, route.Name)

		if route.Securitize {
			handler = c.JwtMiddleware.Handler(handler)
			//TODO Handle jwt token expired
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"SignUp",
		"POST",
		"/api/dentist/signup",
		d.SignUpRegister,
		false,
	},

	Route{
		"SignUpActivate",
		"GET",
		"/api/dentist/activate",
		d.SignUpActivate,
		false,
	},

	Route{
		"ResetPasswordCode",
		"POST",
		"/api/dentist/sendPasswordResetConfirmationCode",
		d.SendPasswordResetConfirmationCode,
		false,
	},

	Route{
		"ResetPassword",
		"POST",
		"/api/dentist/resetPassword",
		d.ResetPassword,
		false,
	},

	Route{
		"Get",
		"GET",
		"/api/dentist/get",
		d.Get,
		false,
	},

	Route{
		"Login",
		"POST",
		"/api/dentist/login",
		d.Login,
		false,
	},
	Route{
		"SeedData",
		"GET",
		"/api/seedData",
		SeedData,
		false,
	},
	Route{
		"Patients",
		"GET",
		"/api/patients",
		d.GetPatients,
		true,
	},
	Route{
		"PatientUpdate",
		"POST",
		"/api/patients/update",
		d.UpdatePatientProfile,
		true,
	},
	Route{
		"PatientCreate",
		"POST",
		"/api/patients/create",
		d.CreatePatientProfile,
		true,
	},
	Route{
		"PatientRemove",
		"DELETE",
		"/api/patients/remove",
		d.RemovePatientProfile,
		true,
	},
	Route{
		"SeedTeethData",
		"GET",
		"/api/toothStatus/seedTeethData",
		d.GetTeethData,
		true,
	},
	// Manipulations
	Route{
		"AddToothManipulation",
		"POST",
		"/api/toothStatus/addToothManipulation",
		d.AddToothManipulation,
		true,
	},
	Route{
		"RemoveToothManipulation",
		"POST",
		"/api/toothStatus/removeToothManipulation",
		d.RemoveToothManipulation,
		true,
	},
	// Diagnosis
	Route{
		"AddToothDiagnosis",
		"POST",
		"/api/toothStatus/addToothDiagnosis",
		d.AddToothDiagnosis,
		true,
	},
	Route{
		"RemoveToothDiagnosis",
		"POST",
		"/api/toothStatus/removeToothDiagnosis",
		d.RemoveToothDiagnosis,
		true,
	},

	// invitation
	Route{
		"InvitePatient",
		"POST",
		"/api/dentist/invite",
		d.InvitePatient,
		true,
	},
	Route{
		"Invitationactivate",
		"GET",
		"/api/dentist/invitationActivate",
		d.InvitationActivate,
		false,
	},

	// patient
	Route{
		"SignUpPatient",
		"POST",
		"/api/patient/signup",
		p.SignUpRegister,
		false,
	},
	Route{
		"SignUpPatirntActivate",
		"GET",
		"/api/patient/activate",
		p.SignUpActivate,
		false,
	},
	Route{
		"ResetPatientPasswordCode",
		"POST",
		"/api/patient/sendPasswordResetConfirmationCode",
		p.SendPasswordResetConfirmationCode,
		false,
	},
	Route{
		"ResetPatientPassword",
		"POST",
		"/api/patient/resetPassword",
		p.ResetPassword,
		false,
	},
	Route{
		"LoginPatient",
		"POST",
		"/api/patient/login",
		p.Login,
		false,
	},
}
