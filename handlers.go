package main

import (
	"bytes"
	"fmt"
	"html/template"
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

	t, err := template.New("filename").Parse(config.FilenameTemplate)
	if err != nil {
		fmt.Println("1")
		panic(err)
	}

	fmt.Println("2")
	var n bytes.Buffer
	err = t.Execute(&n, p)
	if err != nil {
		fmt.Println("3")
		panic(err)
	}
	fmt.Println("4")
	name := n.String()

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
