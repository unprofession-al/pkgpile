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
		l.l("creating repo", "creating repo "+reponame+" which did not exist")
	}
	store[reponame][sumString] = *p
	l.l("storing package", "package "+n.String()+" is saved")

	repodata[reponame], err = yum.CreateRepoData(store[reponame])
	if err != nil {
		panic(err)
	}
	l.l("updating repodata", "repodata saved")

	r.JSON(res, http.StatusOK, n.String())
}

func GetConfig(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, config)
}

func GetFilelists(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/gzip")
	res.Write(repodata[reponame].Filelists)
}

func GetOther(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/gzip")
	res.Write(repodata[reponame].Other)
}

func GetPrimary(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/gzip")
	res.Write(repodata[reponame].Primary)
}

func GetRepomd(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	res.Write(repodata[reponame].Repomd)
}
