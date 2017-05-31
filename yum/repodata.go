package yum

import "github.com/cavaliercoder/go-rpm"

type RepoData struct {
	Primary   []byte
	Filelists []byte
	Other     []byte
	Repomd    []byte
}

func CreateRepoData(packages map[string]rpm.PackageFile) (RepoData, error) {
	rd := RepoData{}
	req := make(map[string]RepomdRequirements)

	pdata := GetPrimary(packages)
	pstr, pstrsize, pstrsum, err := GetXML(pdata)
	if err != nil {
		return rd, err
	}
	pzip, pzipsize, pzipsum := GetZip(pstr)
	rd.Primary = pzip
	req["primary"] = RepomdRequirements{
		OpenSum:  pstrsum,
		OpenSize: pstrsize,
		Sum:      pzipsum,
		Size:     pzipsize,
	}

	fdata := GetFilelists(packages)
	fstr, fstrsize, fstrsum, err := GetXML(fdata)
	if err != nil {
		return rd, err
	}
	fzip, fzipsize, fzipsum := GetZip(fstr)
	rd.Filelists = fzip
	req["filelists"] = RepomdRequirements{
		OpenSum:  fstrsum,
		OpenSize: fstrsize,
		Sum:      fzipsum,
		Size:     fzipsize,
	}

	odata := GetOther(packages)
	ostr, ostrsize, ostrsum, err := GetXML(odata)
	if err != nil {
		return rd, err
	}
	ozip, ozipsize, ozipsum := GetZip(ostr)
	rd.Primary = ozip
	req["other"] = RepomdRequirements{
		OpenSum:  ostrsum,
		OpenSize: ostrsize,
		Sum:      ozipsum,
		Size:     ozipsize,
	}

	rdata := GetRepomd(req)
	rstr, _, _, err := GetXML(rdata)
	if err != nil {
		return rd, err
	}
	rd.Repomd = rstr

	return rd, nil
}
