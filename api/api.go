package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start(bind string) error {

	r := mux.NewRouter()
	r.Methods("GET").Path("/api/stats").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	return http.ListenAndServe(bind, r)

}
