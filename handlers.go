package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cavaliercoder/go-rpm"
	"github.com/gorilla/mux"
	"github.com/unprofession-al/pkgpile/yum"
	"github.com/unrolled/render"
)

func UploadPackage(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	r := render.New()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	sum := sha256.Sum256(body)
	sumString := fmt.Sprintf("%x", sum)

	p, err := rpm.ReadPackageFile(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	var n bytes.Buffer
	err = filenameTemplate.Execute(&n, p)
	if err != nil {
		panic(err)
	}

	if _, ok := store[reponame]; !ok {
		store[reponame] = map[string]rpm.PackageFile{}
	}
	store[reponame][sumString] = *p

	r.JSON(res, http.StatusOK, n.String())
}

func GetConfig(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, config)
}

func GetFilelists(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	r := render.New(render.Options{
		PrefixXML: []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"),
		IndentXML: true,
	})

	filelists := yum.FilelistsRender(store[reponame])
	r.XML(res, http.StatusOK, filelists)
}
