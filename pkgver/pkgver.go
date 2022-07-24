// Package pkgver implements parsing of xbps pkgver and pkgpatterns.
//
// The parsing follows:
//  pkgver  ::= name [ ( version / pattern ) ]
//  name    ::= [a-zA-Z0-9-]*
//  version ::= "-" [^-_]* [ "_" [0-9]+ ]
//  pattern ::= ( ">=" ">" / "<=" / "<" / "==" / "!=" ) .+
//
// Note that a malformed version without "_" will not result in an error.
// The malformed version part will be part of the name.
package pkgver

import (
	"errors"
	"fmt"
	"strings"
)

// A PkgVer represents a name and optionally version or pattern.
type PkgVer struct {
	Name    string
	Version string
	Pattern string
}

var errPattern = errors.New("malformed pattern")

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

// onlyDigits returns true if all runes in str are digits, otherwise false
func onlyDigits(str string) bool {
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// duckPkgver looks if the string following the last - looks like a version.
func duckPkgver(s string) PkgVer {
	len := len(s)
	if i := strings.LastIndexByte(s, '-'); i != -1 {
		// if it contains _, it quacks like a version
		rev := strings.LastIndexByte(s[i+1:], '_')
		if i == len || rev == -1 || i+1+rev+1 == len || !onlyDigits(s[i+1+rev+1:]) {
			return PkgVer{Name: s}
		}
		return PkgVer{Name: s[:i], Version: s[i+1:]}
	}
	return PkgVer{Name: s}
}

// parsePattern checks if the pattern is followed by at least one character
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
