package yum

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Repomd struct {
	XMLName  xml.Name     `xml:"repomd"`
	Revision int          `xml:"revision"`
	Data     []RepomdData `xml:"data"`
	Xmlns    string       `xml:"xmlns,attr"`
	XmlnsRpm string       `xml:"xmlns:rpm,attr"`
}

type RepomdData struct {
	Checksum     Checksum `xml:"checksum"`
	OpenChecksum Checksum `xml:"open-checksum"`
	Location     Location `xml:"location"`
	Type         string   `xml:"type,attr"`
	Timestamp    int      `xml:"timestamp"`
	Size         int      `xml:"size"`
	OpenSize     int      `xml:"open-size"`
}

type Checksum struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Location struct {
	Href string `xml:"href,attr"`
}

func RepomdGen() {
	data := []RepomdData{
		RepomdData{
			Type: "filelists",
			Checksum: Checksum{
				Type:  "sha256",
				Value: "b4db538247fe3cc2c6e97c7f3abed9a82b62512b2d4db3ce1185c5eb6eb5dde8",
			},
			OpenChecksum: Checksum{
				Type:  "sha256",
				Value: "35731cb3f66e5c0930532d00bec62abced097b6ab79b3deb597b77544bf9eddc",
			},
			Location: Location{
				Href: "repodata/b4db538247fe3cc2c6e97c7f3abed9a82b62512b2d4db3ce1185c5eb6eb5dde8-filelists.xml.gz",
			},
			Timestamp: 1495185170,
			Size:      1752970,
			OpenSize:  28290543,
		},
	}
	r := Repomd{
		Revision: 1495185162,
		Xmlns:    "http://linux.duke.edu/metadata/repo",
		XmlnsRpm: "http://linux.duke.edu/metadata/rpm",
		Data:     data,
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(&r); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
