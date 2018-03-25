package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
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

	return http.ListenAndServe(bind, r)

}
