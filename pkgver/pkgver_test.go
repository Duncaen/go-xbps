package pkgver

import (
	"testing"
)

var parseTests = []struct {
	str string
	res PkgVer
}{
	{"foo-1.0_1", PkgVer{Name: "foo", Version: "1.0_1"}},
	{"foo", PkgVer{Name: "foo"}},
	{"foo-32bit-1.0_1", PkgVer{Name: "foo-32bit", Version: "1.0_1"}},
	{"foo-32bit", PkgVer{Name: "foo-32bit"}},
	{"foo-32bit-1.0", PkgVer{Name: "foo-32bit-1.0"}},

	{"foo>1.0_1", PkgVer{Name: "foo", Pattern: ">1.0_1"}},
	{"foo-32bit>=1.0_1", PkgVer{Name: "foo-32bit", Pattern: ">=1.0_1"}},
	{"foo-32bit==1.0_1", PkgVer{Name: "foo-32bit", Pattern: "==1.0_1"}},
	{"foo-32bit!=1.0", PkgVer{Name: "foo-32bit", Pattern: "!=1.0"}},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		pkgver := Parse(tt.str)
		if pkgver != tt.res {
			t.Fatalf("expected %v, got %v", tt.res, pkgver)
		}
		if pkgver.String() != tt.str {
			t.Fatalf("expected string representation %q, got %q", tt.res, pkgver)
		}
	}
}
