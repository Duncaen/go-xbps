package pkgver

import (
	"errors"
	"fmt"
	"strings"
)

// A PkgVer represents a name and a version or pattern.
//
// The format follows:
//  pkgver  ::= name (version / pattern)
//  name    ::= [a-zA-Z0-9-]*
//  version ::= "-" [^-_]* "_" [^-]*
//  pattern ::= ("=>" ">" / "<=" / "<" / "==" / "!=") .*
//
// Note that a malformed version without "_" will not result in an error.
// The malformed version part will be part of the name.
type PkgVer struct {
	Name    string
	Version string
	Pattern string
}

var (
	errPattern = errors.New("malformed pattern")
)

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

func parsePattern(s string, i int) (PkgVer, error) {
	c, l := s[i], len(s[i:])
	if ((c == '!' || c == '=') && (l < 3 || s[i+1] != '=')) ||
		(c == '>' || c == '<') && (l < 2 || (s[i+1] == '=' && l < 3)) {
		return PkgVer{Name: s}, errPattern
	}
	return PkgVer{Name: s[:i], Pattern: s[i:]}, nil
}

// Parse splits package name strings into name, version and pattern parts.
func Parse(s string) (PkgVer, error) {
	if i := strings.IndexAny(s, "><=!"); i != -1 {
		return parsePattern(s, i)
	}
	return duckPkgver(s), nil
}
