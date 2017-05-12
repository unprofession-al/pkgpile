package yum

import (
	"errors"
	"strings"
)

type Rpm struct {
	Name string
	Arch string
}

func NewRpm(filename string) (*Rpm, error) {
	n, a, err := getInfo(filename)
	if err != nil {
		return nil, err
	}

	rpm := &Rpm{
		Name: n,
		Arch: a,
	}

	return rpm, nil
}

func getInfo(filename string) (name string, arch string, err error) {
	err = nil

	tokens := strings.Split(filename, ".")
	if len(tokens) < 3 {
		err = errors.New("Package name must consist of [name].[arch].rpm, eg. 'rtmpdump-2.3-1.el7.x86_64.rpm'")
		return
	}
	if tokens[len(tokens)-1] != "rpm" {
		err = errors.New("Not a RPM package")
		return
	}

	arch = tokens[len(tokens)-2]
	name = strings.Join(tokens[0:len(tokens)-3], ".")

	return
}
