package yum

import "encoding/xml"

type Filelists struct {
	XMLName  xml.Name           `xml:"filelists"`
	Packages int                `xml:"packages"`
	Package  []FilelistsPackage `xml:"package"`
	Xmlns    string             `xml:"xmlns,attr"`
}

type FilelistsPackage struct {
	Version      FilelistsVersion `xml:version`
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
	Version int    `xml:"ver,attr"`
	Release string `xml:"rel,attr"`
}
