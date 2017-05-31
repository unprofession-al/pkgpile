package yum

import "github.com/cavaliercoder/go-rpm"

type RepoData map[string][]byte

func CreateRepoData(packages map[string]rpm.PackageFile) (RepoData, error) {
	rd := RepoData{}
	req := make(map[string]RepomdRequirements)

	pdata := GetPrimary(packages)
	pstr, pstrsize, pstrsum, err := GetXML(pdata)
	if err != nil {
		return rd, err
	}
	pzip, pzipsize, pzipsum := GetZip(pstr)
	rd["primary.xml"] = pstr
	rd["primary.xml.gz"] = pzip
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
	rd["filelists.xml"] = fstr
	rd["filelists.xml.gz"] = fzip
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
	rd["other.xml"] = ostr
	rd["other.xml.gz"] = ozip
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
	rzip, _, _ := GetZip(rstr)
	rd["repomd.xml"] = rstr
	rd["repomd.xml.gz"] = rzip

	return rd, nil
}
