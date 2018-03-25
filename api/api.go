package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/netice9/swarm-intelligence/core"
)

func Start(bind string) error {

	r := mux.NewRouter()
	r.Methods("POST").Path("/api/deploy_stack").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(0)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		f, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		tf, err := ioutil.TempFile("", "stack")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		defer func() {
			tf.Close()
			os.Remove(tf.Name())
		}()

		_, err = io.Copy(tf, f)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		cmd := exec.Command("docker", "stack", "deploy", "-c", tf.Name(), "--prune", "--with-registry-auth", r.FormValue("name"))
		buf := &bytes.Buffer{}
		cmd.Stdout = buf
		cmd.Stderr = buf
		err = cmd.Run()

		cmd.StdoutPipe()

		if err != nil {
			http.Error(w, buf.String(), 500)
		}

	})

	r.Methods("GET").Path("/api/state").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Println(err.Error())
			return
		}
		for {
			s := core.CurrentState()
			err = c.WriteJSON(s)
			if err != nil {
				return
			}
			time.Sleep(2 * time.Second)
		}
	})

	return http.ListenAndServe(bind, r)

}
