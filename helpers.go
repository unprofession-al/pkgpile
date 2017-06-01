package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/unprofession-al/pkgpile/yum"
)

func getBodyAsBytes(body io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createRepo(name string) error {
	if _, ok := metadata[name]; !ok {
		metadata[name] = yum.PackageInfos{}
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
