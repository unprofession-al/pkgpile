package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

	pi := yum.PackageInfo{n.String(), *p}

	if _, ok := metadata[reponame]; !ok {
		metadata[reponame] = yum.PackageInfos{}
		l.l("creating repo", "creating repo "+reponame+" which did not exist")
	}
	metadata[reponame][sumString] = pi
	l.l("storing package", "package "+n.String()+" is saved")

	repodata[reponame], err = yum.CreateRepoData(metadata[reponame])
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

func GetRepoData(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]
	file := vars["file"]
	if data, ok := repodata[reponame][file]; ok {
		res.WriteHeader(http.StatusOK)
		if strings.HasSuffix(file, "gz") {
			res.Header().Set("Content-Type", "application/gzip")
		} else {
			res.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		}
		res.Write(data)
		return
	}
	res.WriteHeader(http.StatusNotFound)
}
