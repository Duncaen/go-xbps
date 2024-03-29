package repo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

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
	Arch     string
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

func New(url, arch string) (*Repository, error) {
	uri, err := uri.Parse(url)
	if err != nil {
		return nil, err
	}
	return &Repository{URI: uri, Arch: arch}, nil
}

func (r *Repository) Sync() error {
	return errors.New("not implemented")
}

func Open(url, arch string) (*Repository, error) {
	uri, err := uri.Parse(url)
	if err != nil {
		return nil, err
	}
	r := &Repository{URI: uri, Arch: arch}
	err = r.Open()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repository) Open() error {
	var repodata, stagedata string
	switch r.URI.Scheme {
	case "file", "":
		repodata = path.Join(r.URI.Path, fmt.Sprintf("%s-repodata", r.Arch))
		stagedata = path.Join(r.URI.Path, fmt.Sprintf("%s-stagedata", r.Arch))
	default:
		return errors.New("not implemented")
	}
	f, err := os.Open(repodata)
	if err != nil {
		return err
	}
	repopkgs, err := r.Read(f)
	if err != nil {
		f.Close()
		return err
	}
	r.Packages = repopkgs
	f.Close()
	f, err = os.Open(stagedata)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	stagepkgs, err := r.Read(f)
	if err != nil {
		f.Close()
		return err
	}
	for k, v := range stagepkgs {
		r.Packages[k] = v
	}
	return nil
}

func (r *Repository) Read(f io.Reader) (map[string]Package, error) {
	rd, err := NewReader(f)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	var packages map[string]Package
	for {
		name, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch name {
		case indexFile:
			packages, err = rd.ReadPackages()
			if err != nil {
				return nil, err
			}
		}
	}
	return packages, nil
}
