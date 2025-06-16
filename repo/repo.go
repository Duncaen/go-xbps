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
	SourcePkg       string              `plist:"sourcepkg"`
}

type Repository struct {
	Arch           string
	URI            *uri.URI
	Packages       map[string]Package
	StagedPackages map[string]Package
}

// New create a new repository structure
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

// Open opens and reads a new repository
func Open(url, arch string) (*Repository, error) {
	uri, err := uri.Parse(url)
	if err != nil {
		return nil, err
	}
	r := &Repository{URI: uri, Arch: arch}
	if err := r.Open(); err != nil {
		return nil, err
	}
	return r, nil
}

// Open reads the repository data from the repositories uri
func (r *Repository) Open() error {
	var repodata string
	switch r.URI.Scheme {
	case "file", "":
		repodata = path.Join(r.URI.Path, fmt.Sprintf("%s-repodata", r.Arch))
	default:
		return fmt.Errorf("repo scheme not supported: %s", r.URI.Scheme)
	}
	f, err := os.Open(repodata)
	if err != nil {
		return fmt.Errorf("repo could not be opened: %w", err)
	}
	defer f.Close()
	if _, err := r.ReadFrom(f); err != nil {
		return fmt.Errorf("repo could not be read: %w", err)
	}
	return nil
}

// ReadFrom reads the repository data from the reader
func (r *Repository) ReadFrom(rd io.Reader) (int64, error) {
	dec, err := NewDecoder(rd)
	if err != nil {
		return 0, err
	}
	defer dec.Close()

	for {
		name, err := dec.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return dec.reader.n, fmt.Errorf("failed to read repository: read header: %w", err)
		}
		switch name {
		case IndexEntry:
			if err := dec.ReadPlist(&r.Packages); err != nil {
				return dec.reader.n, fmt.Errorf("failed to read repository: read packages: %w", err)
			}
		case StageEntry:
			if err := dec.ReadPlist(&r.StagedPackages); err != nil {
				return dec.reader.n, fmt.Errorf("failed to read repository: read staged packages: %w", err)
			}
		}
	}
	return dec.reader.n, nil
}
