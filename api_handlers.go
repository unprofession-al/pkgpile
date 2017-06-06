package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func GetConfig(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(res, http.StatusOK, config)
}

func ListRepos(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	out := []string{}
	for repo, _ := range packageinfo {
		out = append(out, repo)
	}
	r.JSON(res, http.StatusOK, out)
}

func ListPackages(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)
	reponame := vars["repo"]
	out := []string{}
	if pkgs, ok := packageinfo[reponame]; ok {
		for _, pkg := range pkgs {
			out = append(out, pkg.Path)
		}
		r.JSON(res, http.StatusOK, out)
		return
	}
	r.JSON(res, http.StatusNotFound, "Repo does not exist")
}

func GetPackageInfo(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)
	reponame := vars["repo"]
	pkgname := vars["package"]
	if repo, ok := packageinfo[reponame]; ok {
		for _, pkg := range repo {
			if pkg.Path == pkgname {
				r.JSON(res, http.StatusOK, pkg)
				return
			}
		}
		r.JSON(res, http.StatusNotFound, "Package does not exist")
		return
	}
	r.JSON(res, http.StatusNotFound, "Repo does not exist")
}
