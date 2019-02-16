package pkgver

import (
	"fmt"
	"strings"
)

type PkgVer struct {
	Name    string
	Version string
}

func (p PkgVer) String() string {
	if p.Version != "" {
		return fmt.Sprintf("%s-%s", p.Name, p.Version)
	}
	return p.Name
}

func duckPkgver(s string) PkgVer {
	if i := strings.LastIndexByte(s, '-'); i != -1 {
		// if it contains _, it quacks like a version
		if i == len(s) || strings.IndexByte(s[i+1:], '_') == -1 {
			return PkgVer{Name: s}
		}
		return PkgVer{Name: s[:i], Version: s[i+1:]}
	}
	return PkgVer{Name: s}
}

func Parse(s string) PkgVer {
	if i := strings.IndexAny(s, "><=!"); i != -1 {
		return PkgVer{Name: s[:i]}
	}
	return duckPkgver(s)
}
