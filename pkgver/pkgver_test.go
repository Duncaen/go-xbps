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

	{nil, "perl-Text-CSV_XS", PkgVer{Name: "perl-Text-CSV_XS"}},
	{nil, "perl-Text-CSV_", PkgVer{Name: "perl-Text-CSV_"}},
	{nil, "perl-Text-", PkgVer{Name: "perl-Text-"}},
	{nil, "perl-Text-_", PkgVer{Name: "perl-Text-_"}},

	{nil, "perl-Text-CSV_XS-1.40_1", PkgVer{Name: "perl-Text-CSV_XS", Version: "1.40_1"}},

	{nil, "perl-Digest-1.17_01_1", PkgVer{Name: "perl-Digest", Version: "1.17_01_1"}},
	{nil, "perl-PerlIO-utf8_strict-0.007_1", PkgVer{Name: "perl-PerlIO-utf8_strict", Version: "0.007_1"}},

	{nil, "perl-PerlIO-utf8_strict", PkgVer{Name: "perl-PerlIO-utf8_strict" }},
	{nil, "font-adobe-100dpi-1.8_blah", PkgVer{Name: "font-adobe-100dpi-1.8_blah"}},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		pkgver, err := Parse(tt.str)
		if err != tt.err {
			t.Fatalf("expected error %v, got %v", tt.err, err)
		}
		if pkgver != tt.res {
			t.Fatalf("expected %#v, got %#v", tt.res, pkgver)
		}
		if pkgver.String() != tt.str {
			t.Fatalf("expected string representation %q, got %q", tt.res, pkgver)
		}
	}
}
