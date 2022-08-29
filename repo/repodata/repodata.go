package repodata

import (
	"io"
)

const (
	IndexFile     = "index.plist"
	IndexMetaFile = "index-meta.plist"
)

type Package struct {
	Alternatives    map[string][]string `plist:"alternatives"`
	Architecture    string              `plist:"architecture"`
	BuildDate       string              `plist:"build-date"`
	BuildOptions    string              `plist:"build-options"`
	ConfFiles       []string            `plist:"conf_files"`
	Conflicts       []string            `plist:"conflicts"`
	FilenameSHA256  string              `plist:"filename-sha256"`
	FilenameSize    int64               `plist:"filename-size"`
	Homepage        string              `plist:"homepage"`
	InstalledSize   int64               `plist:"installed_size"`
	License         string              `plist:"license"`
	Maintainer      string              `plist:"maintainer"`
	PkgVer          string              `plist:"pkgver"`
	Preserve        bool                `plist:"preserve"`
	Replaces        []string            `plist:"replaces"`
	Reverts         []string            `plist:"reverts"`
	RunDepends      []string            `plist:"run_depends"`
	ShlibProvides   []string            `plist:"shlib-provides"`
	ShlibRequires   []string            `plist:"shlib-requires"`
	ShortDesc       string              `plist:"short_desc"`
	SourceRevisions string              `plist:"source-revisions"`
}

type Meta struct {
	Key      []byte `plist:"public-key"`
	Size     uint16 `plist:"public-key-size"`
	SignedBy string `plist:"signature-by"`
}

type RepoData struct {
	Index map[string]Package `repodata:"index.plist"`
	Meta  Meta               `repodata:"index-meta.plist"`
}

// Read creates Data by reading packages and public key from rd
func Read(rd io.ReadSeeker) (*RepoData, error) {
	repodata := &RepoData{}
	err := NewDecoder(rd).Decode(repodata)
	if err != nil {
		return nil, err
	}
	return repodata, nil
}
