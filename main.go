package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sontags/env"
	"github.com/unprofession-al/pkgpile/yum"
)

var version string
var commitId string

var config = &Configuration{
	Version:  version,
	CommitId: commitId,
}

var l = NewLogger()
var filenameTemplate *template.Template

func init() {
	env.Var(&config.Port, "PORT", "8080", "Port to bind")
	env.Var(&config.Base, "BASE", "/tmp/pkgpile", "Base directories of the repos")
	env.Var(&config.DebugStr, "DEBUG", "false", "Turn debugging on (only print commands to be run)")
	env.Var(&config.FilenameTemplate, "FILENAME_TEMPLATE", "{{.Name}}-{{.Version}}-{{.Release}}.{{.Architecture}}.rpm", "Turn debugging on (only print commands to be run)")
}

var metadata = map[string]yum.PackageInfos{}
var repodata = map[string]yum.RepoData{}

func main() {
	env.Parse("PKGPILE", false)

	var err error
	filenameTemplate, err = template.New("filename").Parse(config.FilenameTemplate)
	if err != nil {
		panic(err)
	}

	config.Debug = !strings.Contains(config.DebugStr, "false")

	err = readRepos()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{repo}/", UploadPackage).Methods("POST")
	r.HandleFunc("/{repo}/repodata/{file}", GetRepoData).Methods("GET")
	r.HandleFunc("/{repo}/{file}", GetPackage).Methods("GET")
	r.HandleFunc("/config.json", GetConfig).Methods("GET")
	chain := alice.New().Then(r)

	l.l("starting...", "pkgpile should be ready")

	log.Fatal(http.ListenAndServe(":"+config.Port, chain))
}
