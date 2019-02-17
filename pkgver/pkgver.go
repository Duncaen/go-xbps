package pkgver

import (
	"fmt"
	"strings"
)

// A PkgVer represents a pkgname, pkgver or pkgpattern
type PkgVer struct {
	Name    string
	Version string
	Pattern string
}

func (p PkgVer) String() string {
	switch {
	case p.Version != "":
		return fmt.Sprintf("%s-%s", p.Name, p.Version)
	case p.Pattern != "":
		return fmt.Sprintf("%s%s", p.Name, p.Pattern)
	default:
		return p.Name
	}
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

// Parse splits package name strings into name, version and pattern parts.
func Parse(s string) PkgVer {
	if i := strings.IndexAny(s, "><=!"); i != -1 {
		return PkgVer{Name: s[:i], Pattern: s[i:]}
	}
	return duckPkgver(s)
}
