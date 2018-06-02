package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	auth "github.com/nabeken/negroni-auth"
	"github.com/netice9/swarm-intelligence/aggregator"
	"github.com/netice9/swarm-intelligence/core"
	"github.com/netice9/swarm-intelligence/frontend"
	"github.com/urfave/negroni"
)

func Start(bind, aggregatorBind string) error {

	go listenForAggregator(aggregatorBind)

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
		cmd.Stdout = w
		cmd.Stderr = w
		err = cmd.Run()

		if err != nil {
			http.Error(w, buf.String(), 500)
		}

	})

	r.Methods("POST").Path("/api/add_credentials").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		loginData := struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Registry string `json:"registry"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		cmd := exec.Command("docker", "login", loginData.Registry, "-u", loginData.Username, "-p", loginData.Password)
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
			s := aggregator.State()
			err = c.WriteJSON(s)
			if err != nil {
				return
			}
			time.Sleep(2 * time.Second)
		}
	})

	r.Methods("GET").Path("/api/namespaces").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		as := aggregator.State()
		namespaces := []string{}
		for _, se := range as.Services {
			ns := se.Spec.Labels["com.docker.stack.namespace"]
			if ns != "" {
				i := sort.SearchStrings(namespaces, ns)
				if i == len(namespaces) {
					namespaces = append(namespaces, ns)
				}
				if namespaces[i] != ns {
					namespaces = append(namespaces[:i], append([]string{ns}, namespaces[i:]...)...)
				}
			}
		}
		r.Header.Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(namespaces)
	})

	r.Methods("DELETE").Path("/api/namespaces/{namespace}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ns := mux.Vars(r)["namespace"]
		cmd := exec.Command("docker", "stack", "deploy", "rm", ns)
		buf := &bytes.Buffer{}
		cmd.Stdout = w
		cmd.Stderr = w
		err := cmd.Run()

		if err != nil {
			http.Error(w, buf.String(), 500)
		}
	})

	r.Methods("DELETE").Path("/api/services/{serviceID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceID := mux.Vars(r)["serviceID"]
		err := core.DeleteService(serviceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	r.Methods("GET").Path("/api/services/{serviceID}/logs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceID := mux.Vars(r)["serviceID"]

		w.Header().Set("Content-Type", "text/event-stream")

		rc, err := core.ServiceLogs(serviceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rc.Close()

		fw := flushWriter{w}

		stdcopy.StdCopy(fw, fw, rc)

	})

	n := negroni.New()
	authUsername := os.Getenv("AUTH_USERNAME")
	authPassword := os.Getenv("AUTH_PASSWORD")
	if authUsername != "" && authPassword != "" {
		n.Use(auth.Basic(authUsername, authPassword))
	}
	n.Use(negroni.NewStatic(frontend.AssetFS()))
	n.UseHandler(r)
	return http.ListenAndServe(bind, n)

}
