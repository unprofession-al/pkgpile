package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
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
	skipupdate := vars["skipupdate"]

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

	packageinfo[reponame][sumString] = pi
	l.l("storing package metadata", "package "+n.String()+" is saved")

	if skipupdate != "" {
		repodata[reponame], err = yum.CreateRepoData(packageinfo[reponame])
		if err != nil {
			r.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}
		l.l("updating repodata", "repodata saved")
	}

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

func GetRepofile(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reponame := vars["repo"]

	proto := "https"
	if req.TLS == nil {
		proto = "http"
	}

	data := struct {
		Reponame string
		URL      string
	}{
		Reponame: reponame,
		URL:      proto + "://" + req.Host + "/" + reponame + "/",
	}

	t := `[{{.Reponame}}]
name={{.Reponame}}
baseurl={{.URL}}
enabled=1
gpgcheck=0
priority=1`

	templ, err := template.New("repofile").Parse(t)
	if err != nil {
		http.Error(res, "Could not create template.", http.StatusInternalServerError)
		return
	}

	var n bytes.Buffer
	err = templ.Execute(&n, data)
	if err != nil {
		http.Error(res, "Could not render template.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, n.String())
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

func GetRepoDataIndex(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)
	reponame := vars["repo"]
	out := []string{}
	if data, ok := repodata[reponame]; ok {
		for filename, _ := range data {
			out = append(out, filename)
		}
		r.JSON(res, http.StatusOK, out)
		return
	}
	r.JSON(res, http.StatusNotFound, "Repo does not exist")
}

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
				data := struct {
					Path        string `json:"path"`
					Name        string `json:"name"`
					Version     string `json:"version"`
					Release     string `json:"release"`
					Arch        string `json:"arch"`
					Summary     string `json:"summary"`
					Description string `json:"description"`
				}{
					Path:        pkg.Path,
					Name:        pkg.Name(),
					Version:     pkg.Version(),
					Release:     pkg.Release(),
					Arch:        pkg.Architecture(),
					Summary:     pkg.Summary(),
					Description: pkg.Description(),
				}
				r.JSON(res, http.StatusOK, data)
				return
			}
		}
		r.JSON(res, http.StatusNotFound, "Package does not exist")
		return
	}
	r.JSON(res, http.StatusNotFound, "Repo does not exist")
}
