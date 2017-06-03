package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/go-rpm"
	"github.com/unprofession-al/pkgpile/yum"
)

func getBodyAsBytes(body io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func createRepo(name string) error {
	if _, ok := packageinfo[name]; !ok {
		packageinfo[name] = yum.PackageInfos{}
		l.l("creating repo data store", "created repo datastore for "+name+" which did not exist")
	}
	path := config.Base + "/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		l.l("creating repo directory on filesystem", "created repo filesystem for "+name+" which did not exist")
		return err
	}
	return nil
}

func createBaseDir() error {
	if _, err := os.Stat(config.Base); os.IsNotExist(err) {
		err := os.MkdirAll(config.Base, 0755)
		l.l("creating base directory on filesystem", "created basedir filesystem: "+config.Base)
		return err
	}
	return nil

}

func reindexPackages() error {
	// reset
	packageinfo = map[string]yum.PackageInfos{}
	repodata = map[string]yum.RepoData{}

	// find repos
	elems, err := ioutil.ReadDir(config.Base)
	if err != nil {
		return err
	}
	for _, elem := range elems {
		if elem.IsDir() {
			l.l("existing repo detected", "exitsing repo "+elem.Name())
			packageinfo[elem.Name()] = yum.PackageInfos{}
		}
	}

	// read rpms in repo
	for repo, _ := range packageinfo {
		l.l("scanning repo", "repo "+repo)
		err = filepath.Walk(config.Base+"/"+repo, func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				if strings.HasSuffix(f.Name(), "rpm") {
					l.l("found package", "pkg "+f.Name())
					// get sha256
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					defer file.Close()
					hasher := sha256.New()
					if _, err := io.Copy(hasher, file); err != nil {
						return err
					}
					sumString := fmt.Sprintf("%x", hasher.Sum(nil))
					l.l("hashing pkg", "hash "+sumString)
					// get rpm info
					p, err := rpm.OpenPackageFile(path)
					if err != nil {
						return err
					}
					pi := yum.PackageInfo{f.Name(), *p}
					// store
					packageinfo[repo][sumString] = pi
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		repodata[repo], err = yum.CreateRepoData(packageinfo[repo])
		if err != nil {
			return err
		}
	}
	return nil
}

func savePackage(data io.Reader, repo string, pkg string) error {
	filepath := config.Base + "/" + repo + "/" + pkg
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	return err
}

func readRepos() error {
	return nil
}
