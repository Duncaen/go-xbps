package pkgver

import (
	"testing"
)

var parseTests = []struct {
	err error
	str string
	res PkgVer
}{
	{nil, "foo-1.0_1", PkgVer{Name: "foo", Version: "1.0_1"}},
	{nil, "foo", PkgVer{Name: "foo"}},
	{nil, "foo-32bit-1.0_1", PkgVer{Name: "foo-32bit", Version: "1.0_1"}},
	{nil, "foo-32bit", PkgVer{Name: "foo-32bit"}},
	{nil, "foo-32bit-1.0", PkgVer{Name: "foo-32bit-1.0"}},

	{nil, "foo>1.0_1", PkgVer{Name: "foo", Pattern: ">1.0_1"}},
	{nil, "foo-32bit>=1.0_1", PkgVer{Name: "foo-32bit", Pattern: ">=1.0_1"}},
	{nil, "foo-32bit==1.0_1", PkgVer{Name: "foo-32bit", Pattern: "==1.0_1"}},
	{nil, "foo-32bit!=1.0", PkgVer{Name: "foo-32bit", Pattern: "!=1.0"}},

	{errPattern, "foo>", PkgVer{Name: "foo>"}},
	{nil, "foo>a", PkgVer{Name: "foo", Pattern: ">a"}},
	{errPattern, "foo>=", PkgVer{Name: "foo>="}},
	{nil, "foo>=a", PkgVer{Name: "foo", Pattern: ">=a"}},
	{errPattern, "foo=", PkgVer{Name: "foo="}},
	{errPattern, "foo==", PkgVer{Name: "foo=="}},
	{nil, "foo==a", PkgVer{Name: "foo", Pattern: "==a"}},
	{errPattern, "foo!", PkgVer{Name: "foo!"}},
	{errPattern, "foo!=", PkgVer{Name: "foo!="}},
	{nil, "foo!=a", PkgVer{Name: "foo", Pattern: "!=a"}},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		pkgver, err := Parse(tt.str)
		if err != tt.err {
			t.Fatalf("expected error %v, got %v", tt.err, err)
		}
		if pkgver != tt.res {
			t.Fatalf("expected %v, got %v", tt.res, pkgver)
		}
		if pkgver.String() != tt.str {
			t.Fatalf("expected string representation %q, got %q", tt.res, pkgver)
		}
	}
}
