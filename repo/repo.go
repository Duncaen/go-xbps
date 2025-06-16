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
	Arch           string
	URI            *uri.URI
	Packages       map[string]Package
	StagedPackages map[string]Package
	Staged         bool
}

const (
	indexFile     = "index.plist"
	stageFile     = "stage.plist"
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
	var repodata string
	switch r.URI.Scheme {
	case "file", "":
		repodata = path.Join(r.URI.Path, fmt.Sprintf("%s-repodata", r.Arch))
	default:
		return errors.New("not implemented")
	}
	f, err := os.Open(repodata)
	if err != nil {
		return err
	}
	err = r.read(f)
	if err != nil {
		f.Close()
		return err
	}
	f.Close()
	return nil
}

func (r *Repository) read(f io.Reader) error {
	rd, err := NewReader(f)
	if err != nil {
		return err
	}
	defer rd.Close()
	var packages map[string]Package
	var stagePackages map[string]Package
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
			packages, err = rd.ReadPackages()
			if err != nil {
				return err
			}
		case stageFile:
			stagePackages, err = rd.ReadPackages()
			if err != nil {
				return err
			}
		}
	}
	if len(stagePackages) > 0 {
		r.Staged = true
	} else {
		r.Staged = false
	}
	r.Packages = packages
	r.StagedPackages = stagePackages
	return nil
}
