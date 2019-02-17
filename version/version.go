// Package version implements the dewey algorithm xbps uses.
// It is based on NetBSDs algorithm, but removed the NetBSD nb
// modifier and added a revision field indicated by underscore (_).
package version

import (
	"strings"
)

type op int

const (
	oLT op = iota
	oLE
	oEQ
	oGE
	oGT
	oNE
)

// Version represents a parsed version string.
type Version struct {
	v []int
	r int
}

type parse struct {
	buf string
	len int
	pos int
	arr []int
	rev int
}

// DO NOT MODIFY these values, or things will NOT work.
const (
	alpha = -3
	beta  = -2
	rc    = -1
	dot   = 0
)

// modifier strings and what version and what version they represent.
var modifiers = []struct {
	s string
	n int
}{
	{"alpha", alpha},
	{"beta", beta},
	{"pre", rc},
	{"rc", rc},
	{"pl", dot},
	{".", dot},
}

func (p *parse) modifier() bool {
	for _, mod := range modifiers {
		if strings.HasPrefix(p.buf[p.pos:], mod.s) {
			p.arr = append(p.arr, mod.n)
			p.pos += len(mod.s)
			return true
		}
	}
	return false
}

func (p *parse) number() int {
	n := 0
	for ; p.pos < p.len; p.pos++ {
		c := p.buf[p.pos]
		if !(c >= '0' && c <= '9') {
			// p.pos -= 1
			break
		}
		n = (n * 10) + (int(c) - '0')
	}
	return n
}

// Parse parses a raw version string
func Parse(s string) Version {
	p := parse{buf: strings.ToLower(s), len: len(s)}
	for p.pos < p.len {
		c := p.buf[p.pos]
		switch {
		case c >= '0' && c <= '9':
			p.arr = append(p.arr, p.number())
		case p.modifier():
		case c == '_':
			p.pos++
			p.rev = p.number()
		case (c >= 'a' && c <= 'z'):
			p.arr = append(p.arr, dot)
			p.arr = append(p.arr, (int(c)-'a')+1)
			p.pos++
		default:
			p.pos++
		}
	}
	return Version{v: p.arr, r: p.rev}
}

func result(cmp int, o op) bool {
	switch o {
	case oLT:
		return cmp < 0
	case oLE:
		return cmp <= 0
	case oGT:
		return cmp > 0
	case oGE:
		return cmp >= 0
	case oEQ:
		return cmp == 0
	case oNE:
		return cmp != 0
	default:
		return false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b

}

func digit(v Version, i int) int {
	if i >= len(v.v) {
		return 0
	}
	return v.v[i]
}

func vcmp(lhs Version, o op, rhs Version) bool {
	for i := 0; i < max(len(lhs.v), len(rhs.v)); i++ {
		if cmp := digit(lhs, i) - digit(rhs, i); cmp != 0 {
			return result(cmp, o)
		}
	}
	return result(lhs.r-rhs.r, o)
}

// Cmp compares a and b and returns: returns an integer comparing two version strings.
//  -1 if a <  b
//   0 if a == b
//  +1 if a >  b
func (a Version) Cmp(b Version) int {
	switch {
	case vcmp(a, oLT, b):
		return -1
	case vcmp(a, oGT, b):
		return 1
	default:
		return 0
	}
}

// Cmp compares a and b and returns: returns an integer comparing two version strings.
//  -1 if a <  b
//   0 if a == b
//  +1 if a >  b
func Cmp(a string, b string) int {
	return Parse(a).Cmp(Parse(b))
}
