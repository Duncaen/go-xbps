// Package repo implements reading xbps repository files.
//
// There are two methods for reading repository files:
//  1. Using the Repository structure and associated functions.
//  2. Using the Decoder to manually decode the repository file
//     which allows to skip over files and package metadata that
//     is not required.
package repo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// Meta is a legacy xbps RSA public key
type Meta struct {
	Key      []byte `plist:"public-key"`
	Size     uint16 `plist:"public-key-size"`
	SignedBy string `plist:"signature-by"`
}

// Repository is the parsed repository file
type Repository struct {
	// Arch is the repository architecture
	Arch string
	// URI is the parsed repository URI
	URI *uri.URI
	// Meta is the repositories legacy xbps RSA public key
	Meta *Meta
	// Index is the repository index, mapping package names to packages
	Index map[string]Package
	// stage is the repository staging index, mapping package names to packages
	Stage map[string]Package
}

func formatpath(u *uri.URI, arch, cachedir string) (string, error) {
	switch u.Scheme {
	case "file", "":
		return filepath.Join(u.Path, fmt.Sprintf("%s-repodata", arch)), nil
	case "http", "https":
		if cachedir != "" {
			return filepath.Join(cachedir, u.CacheString(), fmt.Sprintf("%s-repodata", arch)), nil
		}
		return "", fmt.Errorf("repo scheme not supported without cachedir: %s", u.Scheme)
	default:
		return "", fmt.Errorf("repo scheme not supported: %s", u.Scheme)
	}
}

// Path returns the path to the repodata file for url and arch
func Path(url, arch, cachedir string) (string, error) {
	u, err := uri.Parse(url)
	if err != nil {
		return "", err
	}
	return formatpath(u, arch, cachedir)
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
	repo := &Repository{URI: uri, Arch: arch}
	if err := repo.Open(); err != nil {
		return nil, err
	}
	return repo, nil
}

// Open reads the repository data from the repositories uri
func (repo *Repository) Open() error {
	var repodata string
	repodata, err := formatpath(repo.URI, repo.Arch, "")
	if err != nil {
		return fmt.Errorf("repo could not be openend: %w", err)
	}
	f, err := os.Open(repodata)
	if err != nil {
		return fmt.Errorf("repo could not be opened: %w", err)
	}
	defer f.Close()
	if _, err := repo.ReadFrom(f); err != nil {
		return fmt.Errorf("repo could not be read: %w", err)
	}
	return nil
}

// ReadFrom reads the repository data from the reader
func (repo *Repository) ReadFrom(rd io.Reader) (int64, error) {
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
			if err := dec.ReadPlist(&repo.Index); err != nil {
				return dec.reader.n, fmt.Errorf("failed to read repository: read packages: %w", err)
			}
		case StageEntry:
			if err := dec.ReadPlist(&repo.Stage); err != nil {
				return dec.reader.n, fmt.Errorf("failed to read repository: read staged packages: %w", err)
			}
		case MetaEntry:
			repo.Meta = &Meta{}
			if err := dec.ReadPlist(repo.Meta); err != nil {
				return dec.reader.n, fmt.Errorf("failed to read repository: read : %w", err)
			}
		}
	}
	return dec.reader.n, nil
}
