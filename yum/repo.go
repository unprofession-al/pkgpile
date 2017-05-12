package yum

import (
	"os"
	"os/exec"
)

type Repo struct {
	Location string
}

func NewRepo(l string) *Repo {
	return &Repo{
		Location: l,
	}
}

func (r *Repo) Create() ([]byte, error) {
	cmd := exec.Command("createrepo", "--database", r.Location)
	cmd.Env = append(cmd.Env, os.Environ()...)
	return cmd.Output()
}

func (r *Repo) Update() ([]byte, error) {
	cmd := exec.Command("createrepo", "--update", r.Location)
	cmd.Env = append(cmd.Env, os.Environ()...)
	return cmd.Output()
}
