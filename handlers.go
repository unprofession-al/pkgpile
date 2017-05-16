package main

import (
	"bytes"
	"net/http"

	"github.com/cavaliercoder/go-rpm"
	"github.com/unrolled/render"
)

func UploadPackage(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	p, err := rpm.ReadPackageFile(req.Body)
	if err != nil {
		panic(err)
	}

	var n bytes.Buffer
	err = filenameTemplate.Execute(&n, p)
	if err != nil {
		panic(err)
	}

	r.JSON(res, http.StatusOK, n.String())
}

func CreateRepo(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, "Repo created successfully")
}

func GetConfig(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, config)

}
