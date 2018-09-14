package core

import "net/http"

type AppHandler func(http.ResponseWriter, *http.Request) error

// Interceptor for errors
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
