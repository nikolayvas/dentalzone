package core

import (
	"log"
	"net/http"
)

// AppHandler ..
type AppHandler func(http.ResponseWriter, *http.Request) error

// Interceptor for errors
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			"Exception: ",
			err.Error(),
		)

		http.Error(w, "Internal error" /*err.Error()*/, http.StatusInternalServerError)
	}
}
