package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	sum := sha256.Sum256(body)
	sumString := fmt.Sprintf("%x", sum)

	p, err := rpm.ReadPackageFile(bytes.NewReader(body))
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	var n bytes.Buffer
	err = filenameTemplate.Execute(&n, p)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	pi := yum.PackageInfo{n.String(), *p}
	err = createRepo(reponame)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	err = savePackage(bytes.NewReader(body), reponame, n.String())
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	metadata[reponame][sumString] = pi
	l.l("storing package metadata", "package "+n.String()+" is saved")

	repodata[reponame], err = yum.CreateRepoData(metadata[reponame])
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}
	l.l("updating repodata", "repodata saved")

	r.JSON(res, http.StatusOK, n.String())
}

func GetPackage(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]
	file := vars["file"]

	filename := config.Base + "/" + reponame + "/" + file

	openfile, err := os.Open(filename)
	defer openfile.Close()
	if err != nil {
		http.Error(res, "File not found.", http.StatusNotFound)
		return
	}

	header := make([]byte, 512)
	openfile.Read(header)
	contentType := http.DetectContentType(header)

	stat, err := openfile.Stat()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	size := strconv.FormatInt(stat.Size(), 10)

	res.Header().Set("Content-Disposition", "attachment; filename="+file)
	res.Header().Set("Content-Type", contentType)
	res.Header().Set("Content-Length", size)

	openfile.Seek(0, 0)
	io.Copy(res, openfile)
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
