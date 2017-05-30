package yum

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/cavaliercoder/go-rpm"
)

const primaryXmlns = "http://linux.duke.edu/metadata/common"
const primaryXmlnsRpm = "http://linux.duke.edu/metadata/rpm"

type Primary struct {
	XMLName  xml.Name         `xml:"metadata"`
	Packages int              `xml:"packages,attr"`
	Package  []PrimaryPackage `xml:"package"`
	Xmlns    string           `xml:"xmlns,attr"`
	Xmlnsrpm string           `xml:"xmlns:rpm,attr"`
}

type PrimaryPackage struct {
	Type         string        `xml:"type,attr"`
	Name         string        `xml:"name"`
	Architecture string        `xml:"arch"`
	Version      Version       `xml:"version"`
	Checksum     Checksum      `xml:"checksum"`
	Summary      string        `xml:"summary"`
	Description  string        `xml:"description"`
	Packager     string        `xml:"packager"`
	URL          string        `xml:"url"`
	Time         PrimaryTime   `xml:"time"`
	Size         PrimarySize   `xml:"size"`
	Format       PrimaryFormat `xml:"format"`
}

type PrimaryTime struct {
	File  time.Time `xml:"file,attr"`
	Build time.Time `xml:"build,attr"`
}

type PrimarySize struct {
	Package   uint64 `xml:"package,attr"`
	Installed uint64 `xml:"installed,attr"`
	Archived  uint64 `xml:"archived,attr"`
}

type PrimaryFormat struct {
	License     string                   `xml:"rpm:license"`
	Vendor      string                   `xml:"rpm:vendor"`
	Groups      []string                 `xml:"rpm:group"`
	Buildhost   string                   `xml:"rpm:buildhost"`
	SourceRPM   string                   `xml:"rpm:sourcerpm"`
	HeaderRange PrimaryFormatHeaderRange `xml:"rpm:heander-range"`
	Provides    []PrimaryFormatEntry     `xml:"rpm:provides>rpm:entry"`
	Requires    []PrimaryFormatEntry     `xml:"rpm:requires>rpm:entry"`
	Files       []File                   `xml:"file"`
}

type PrimaryFormatHeaderRange struct {
	Start uint64 `xml:"start,attr"`
	End   uint64 `xml:"end,attr"`
}

type PrimaryFormatEntry struct {
	Name    string `xml:"name,attr"`
	Flags   string `xml:"flags,attr,omitempty"`
	Epoch   string `xml:"epoch,attr,omitempty"`
	Version string `xml:"ver,attr,omitempty"`
	Release string `xml:"rel,attr,omitempty"`
	Pre     string `xml:"pre,attr,omitempty"`
}

func PrimaryRender(packages map[string]rpm.PackageFile) Primary {
	primary := Primary{
		Packages: len(packages),
		Xmlns:    primaryXmlns,
		Xmlnsrpm: primaryXmlnsRpm,
		Package:  []PrimaryPackage{},
	}

	for sum, p := range packages {
		pkgversion := Version{
			Epoch:   p.Epoch(),
			Version: p.Version(),
			Release: p.Release(),
		}
		pkgsum := Checksum{
			Value: sum,
			Type:  "sha256",
			Pkgid: "YES",
		}
		pkgtime := PrimaryTime{
			File:  p.FileTime(),
			Build: p.BuildTime(),
		}
		// TODO: Sizes seem not to work quite well
		pkgsize := PrimarySize{
			Package:   p.FileSize(),
			Installed: p.Size(),
			Archived:  p.ArchiveSize(),
		}
		pkgformatheaderrange := PrimaryFormatHeaderRange{
			Start: p.HeaderStart(),
			End:   p.HeaderEnd(),
		}
		pkgformat := PrimaryFormat{
			License:     p.License(),
			Vendor:      p.Vendor(),
			Groups:      p.Groups(),
			Buildhost:   p.BuildHost(),
			SourceRPM:   p.SourceRPM(),
			HeaderRange: pkgformatheaderrange,
			Provides:    []PrimaryFormatEntry{},
			Requires:    []PrimaryFormatEntry{},
			Files:       []File{},
		}
		for _, p := range p.Provides() {
			provided := PrimaryFormatEntry{
				Name:    p.Name(),
				Epoch:   strconv.Itoa(p.Epoch()),
				Release: p.Release(),
				Version: p.Version(),
			}
			pkgformat.Provides = append(pkgformat.Provides, provided)
		}
		for _, r := range p.Requires() {
			requirement := PrimaryFormatEntry{
				Name: r.Name(),
				Pre:  ReadFlags(r.Flags()),
			}
			pkgformat.Requires = append(pkgformat.Requires, requirement)
		}
		for _, f := range p.Files() {
			file := File{
				Value: f.Name(),
			}
			if f.IsDir() {
				// TODO: The if does not quite work
				file.Type = "dir"
			}
			pkgformat.Files = append(pkgformat.Files, file)
		}
		pkg := PrimaryPackage{
			Name:         p.Name(),
			Architecture: p.Architecture(),
			Version:      pkgversion,
			Checksum:     pkgsum,
			Summary:      p.Summary(),
			Description:  p.Description(),
			Packager:     p.Packager(),
			URL:          p.URL(),
			Time:         pkgtime,
			Size:         pkgsize,
			Format:       pkgformat,
		}
		primary.Package = append(primary.Package, pkg)
	}
	return primary
}

func ReadFlags(f int) string {
	var s string
	switch {
	case rpm.DepFlagLesserOrEqual == (f & rpm.DepFlagLesserOrEqual):
		s = fmt.Sprintf("%s <=", s)

	case rpm.DepFlagLesser == (f & rpm.DepFlagLesser):
		s = fmt.Sprintf("%s <", s)

	case rpm.DepFlagGreaterOrEqual == (f & rpm.DepFlagGreaterOrEqual):
		s = fmt.Sprintf("%s >=", s)

	case rpm.DepFlagGreater == (f & rpm.DepFlagGreater):
		s = fmt.Sprintf("%s >", s)

	case rpm.DepFlagEqual == (f & rpm.DepFlagEqual):
		s = fmt.Sprintf("%s =", s)
	}
	return s
}
