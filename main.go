package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sontags/env"
)

var version string
var commitId string

var config = &Configuration{
	Version:  version,
	CommitId: commitId,
}

var l = NewLogger()

func init() {
	env.Var(&config.RecieverPort, "RECIEVER_PORT", "8080", "Port to bind to in order to accept rpms via POST")
	env.Var(&config.ServerPort, "SERVER_PORT", "8081", "Port to bind to in order to serve GETs")
	env.Var(&config.Base, "BASE", "/tmp/", "Base directories of the repos")
	env.Var(&config.DebugStr, "DEBUG", "false", "Turn debugging on (only print commands to be run)")
	env.Var(&config.FilenameTemplate, "FILENAME_TEMPLATE", "{{.Name}}-{{.Version}}-{{.Architecture}}.rpm", "Turn debugging on (only print commands to be run)")
}

func main() {
	env.Parse("YUMR", false)

	config.Debug = !strings.Contains(config.DebugStr, "false")

	r := mux.NewRouter()
	r.HandleFunc("/{repo}", CreateRepo).Methods("POST")
	r.HandleFunc("/{repo}/{filename}", UploadPackage).Methods("POST")
	r.HandleFunc("/config.json", GetConfig).Methods("GET")
	chain := alice.New().Then(r)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(":"+config.RecieverPort, chain))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.Handle("/", http.FileServer(http.Dir(config.Base)))
		log.Fatal(http.ListenAndServe(":"+config.ServerPort, nil))
		wg.Done()
	}()
	wg.Wait()
}
