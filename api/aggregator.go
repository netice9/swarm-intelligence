package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/netice9/swarm-intelligence/aggregator"
	"github.com/netice9/swarm-intelligence/core"
)

func listenForAggregator(bind string) error {
	r := mux.NewRouter()
	r.Methods("POST").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st := core.State{}
		err := json.NewDecoder(r.Body).Decode(&st)
		if err != nil {
			log.Println(err)
			return
		}
		ra := strings.Split(r.RemoteAddr, ":")
		aggregator.NewState(ra[0], st)
	})
	return http.ListenAndServe(bind, r)
}
