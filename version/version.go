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

var ops = []struct {
	str string
	op  Op
}{
	{"<=", LE},
	{"<", LT},
	{">=", GE},
	{">", GT},
	{"==", EQ},
	{"!=", NE},
}

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

// do not modify these values, or things will NOT work
const (
	Alpha = -3
	Beta  = -2
	RC    = -1
	Dot   = 0
	Patch = 1
)

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

func (l Version) Compare(r Version) int {
	if vcmp(l, LT, r) {
		return -1
	} else if vcmp(l, GT, r) {
		return 1
	}
	return 0
}

func (l Version) OpCompare(o Op, r Version) bool { return vcmp(l, o, r) }
func (l Version) LowerThan(r Version) bool       { return vcmp(l, LT, r) }
func (l Version) LowerEqual(r Version) bool      { return vcmp(l, LE, r) }
func (l Version) GreaterThan(r Version) bool     { return vcmp(l, GT, r) }
func (l Version) GreaterEqual(r Version) bool    { return vcmp(l, GE, r) }
func (l Version) Equal(r Version) bool           { return vcmp(l, EQ, r) }
func (l Version) NotEqual(r Version) bool        { return vcmp(l, NE, r) }

func Compare(lhs string, rhs string) int {
	return Parse(lhs).Compare(Parse(rhs))
}

func OpCompare(lhs string, op Op, rhs string) bool {
	return Parse(lhs).OpCompare(op, Parse(rhs))
}
