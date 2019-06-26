package repo

import (
	"os"
	"io"
	"errors"

	"github.com/Duncaen/go-xbps/repo/uri"
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

type Repository struct {
	URI      *uri.URI
	Packages map[string]Package
}

const (
	indexFile = "index.plist"
	indexMetaFile = "index-meta.plist"
)

var (
	errNoIndex = errors.New("repodata does not contain index.plist")
)

func (r *Repository) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.Read(f)
}

func (r *Repository) Read(f io.Reader) error {
	rd, err := NewReader(f)
	if err != nil {
		return err
	}
	defer rd.Close()
	for {
		name, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		switch name {
		case indexFile:
			r.Packages, err = rd.ReadPackages()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

