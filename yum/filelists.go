package yum

import (
	"encoding/xml"

	"github.com/cavaliercoder/go-rpm"
)

const filelistsXmlns = "http://linux.duke.edu/metadata/filelists"

type Filelists struct {
	XMLName  xml.Name           `xml:"filelists"`
	Packages int                `xml:"packages"`
	Package  []FilelistsPackage `xml:"package"`
	Xmlns    string             `xml:"xmlns,attr"`
}

type FilelistsPackage struct {
	Version      FilelistsVersion `xml:"version"`
	File         []FilelistsFile  `xml:"file"`
	Architecture string           `xml:"arch,attr"`
	Pkgid        string           `xml:"pkgid,attr"`
	Name         string           `xml:"name,attr"`
}

type FilelistsFile struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

type FilelistsVersion struct {
	Epoch   int    `xml:"epoch,attr"`
	Version string `xml:"ver,attr"`
	Release string `xml:"rel,attr"`
}

func FilelistsRender(packages map[string]rpm.PackageFile) Filelists {
	filelists := Filelists{
		Packages: len(packages),
		Xmlns:    filelistsXmlns,
		Package:  []FilelistsPackage{},
	}

	for sum, p := range packages {
		pkgversion := FilelistsVersion{
			Epoch:   p.Epoch(),
			Version: p.Version(),
			Release: p.Release(),
		}
		pkgdata := FilelistsPackage{
			Architecture: p.Architecture(),
			Pkgid:        sum,
			Name:         p.Name(),
			Version:      pkgversion,
			File:         []FilelistsFile{},
		}
		for _, f := range p.Files() {
			file := FilelistsFile{
				Value: f.Name(),
			}
			if f.IsDir() {
				file.Type = "dir"
			}
			pkgdata.File = append(pkgdata.File, file)
		}
		filelists.Package = append(filelists.Package, pkgdata)
	}
	return filelists
}
