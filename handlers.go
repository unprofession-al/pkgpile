package main

import (
	"fmt"
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
	name := fmt.Sprintf("%s-%s-%s.rpm", p.Name(), p.Architecture(), p.Version())

	r.JSON(res, http.StatusOK, name)
}

func CreateRepo(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, "Repo created successfully")
}

func GetConfig(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, config)

}
