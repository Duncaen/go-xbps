package pkgver_test

import (
	"fmt"

	"github.com/Duncaen/go-xbps/pkgver"
)

func ExampleParse() {
	var v pkgver.PkgVer
	v, _ = pkgver.Parse("pkgname")
	fmt.Printf("%#v\n", v)
	v, _ = pkgver.Parse("pkgname-1.0_1")
	fmt.Printf("%#v\n", v)
	v, _ = pkgver.Parse("pkgname>=1.0_1")
	fmt.Printf("%#v\n", v)
	// Output:
	// pkgver.PkgVer{Name:"pkgname", Version:"", Pattern:""}
	// pkgver.PkgVer{Name:"pkgname", Version:"1.0_1", Pattern:""}
	// pkgver.PkgVer{Name:"pkgname", Version:"", Pattern:">=1.0_1"}
}
