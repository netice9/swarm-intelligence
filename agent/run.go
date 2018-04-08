package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/netice9/swarm-intelligence/core"
)

// Run executes the agent loop
func Run(remote string) error {

	for ; ; time.Sleep(2 * time.Second) {
		s := core.CurrentState()

		d, err := json.Marshal(s)
		if err != nil {
			log.Println(err)
			continue
		}

		req, err := http.NewRequest("POST", remote, bytes.NewReader(d))
		if err != nil {
			log.Println(err)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

		resp.Body.Close()
	}

}
