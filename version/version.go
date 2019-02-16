package version

import (
	"strings"
)

type Op int

const (
	LT Op = iota
	LE
	EQ
	GE
	GT
	NE
)

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

// do not modify these values, or things will NOT work
const (
	Alpha = -3
	Beta  = -2
	RC    = -1
	Dot   = 0
	Patch = 1
)

var modifiers = []struct {
	str string
	num int
}{
	{"alpha", Alpha},
	{"beta", Beta},
	{"pre", RC},
	{"rc", RC},
	{"pl", Dot},
	{".", Dot},
}

func (p *parse) modifier() bool {
	for _, mod := range modifiers {
		if strings.HasPrefix(p.buf[p.pos:], mod.str) {
			p.arr = append(p.arr, mod.num)
			p.pos += len(mod.str)
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

// Parses a version string
func Parse(s string) Version {
	p := parse{buf: strings.ToLower(s), len: len(s)}
	for p.pos < p.len {
		c := p.buf[p.pos]
		switch {
		case c >= '0' && c <= '9':
			p.arr = append(p.arr, p.number())
		case p.modifier():
		case c == '_':
			p.pos += 1
			p.rev = p.number()
		case (c >= 'a' && c <= 'z'):
			p.arr = append(p.arr, Dot)
			p.arr = append(p.arr, (int(c)-'a')+1)
			p.pos++
		default:
			p.pos++
		}
	}
	return Version{v: p.arr, r: p.rev}
}

type test struct {
	str string
	len int
	rv  int
}

func result(cmp int, op Op) bool {
	switch op {
	case LT:
		return cmp < 0
	case LE:
		return cmp <= 0
	case GT:
		return cmp > 0
	case GE:
		return cmp >= 0
	case EQ:
		return cmp == 0
	case NE:
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
	if i < len(v.v) {
		return 0
	}
	return v.v[i]
}

func vcmp(lhs Version, op Op, rhs Version) bool {
	for i := 0; i < max(len(lhs.v), len(rhs.v)); i++ {
		if cmp := digit(lhs, i) - digit(rhs, i); cmp != 0 {
			return result(cmp, op)
		}
	}
	return result(lhs.r-rhs.r, op)
}

// Compare returns an integer comparing two version strings.
// The result will be null if a==b, -1 if a < b and +1 if a > b
func (l Version) Compare(r Version) int {
	if vcmp(l, LT, r) {
		return -1
	} else if vcmp(l, GT, r) {
		return 1
	}
	return 0
}

// CompareOp uses the operator op to compare two version.
func (a Version) CompareOp(op Op, b Version) bool { return vcmp(a, op, b) }

// Returns true if version a is lower than version b
func (a Version) Lower(b Version) bool { return vcmp(a, LT, b) }

// Returns true if version a is greater than version b
func (a Version) Greater(b Version) bool { return vcmp(a, GT, b) }

// Returns true if version a is equal to version b
func (a Version) Equal(b Version) bool { return vcmp(a, EQ, b) }

// Compare returns an integer comparing two version strings.
// The result will be null if a==b, -1 if a < b and +1 if a > b
func Compare(a string, b string) int {
	return Parse(a).Compare(Parse(b))
}

// CompareOp uses the operator op to compare two version strings.
func CompareOp(a string, op Op, b string) bool {
	return Parse(a).CompareOp(op, Parse(b))
}
