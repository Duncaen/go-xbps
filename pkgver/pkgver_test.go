package pkgver

import (
	"testing"
)

var parseTests = []struct {
	str string
	res PkgVer
}{
	{"foo-1.0_1", PkgVer{"foo", "1.0_1"}},
	{"foo", PkgVer{"foo", ""}},
	{"foo-32bit-1.0_1", PkgVer{"foo-32bit", "1.0_1"}},
	{"foo-32bit", PkgVer{"foo-32bit", ""}},
	{"foo-32bit-1.0", PkgVer{"foo-32bit-1.0", ""}},

	{"foo>1.0_1", PkgVer{"foo", ""}},
	{"foo-32bit>=1.0_1", PkgVer{"foo-32bit", ""}},
	{"foo-32bit==1.0_1", PkgVer{"foo-32bit", ""}},
	{"foo-32bit!=1.0", PkgVer{"foo-32bit", ""}},
}

func TestPkgVer(t *testing.T) {
	for _, tt := range parseTests {
		pkgver := Parse(tt.str)
		if pkgver != tt.res {
			t.Fatalf("expected %v, got %v", tt.res, pkgver)
		}
	}
}
