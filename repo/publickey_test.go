package repo

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
)

var root = "/"

func TestFilename(t *testing.T) {
	keyfiles, err := filepath.Glob(path.Join(root, "var/db/xbps/keys/*.plist"))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range keyfiles {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var key PublicKey
		if err := ParsePublicKey(buf, &key); err != nil {
			t.Fatal(err)
		}
		res := key.Filename()
		if res != path.Base(f) {
			t.Errorf("Filename() %q does not match %q", res, f)
		}
	}
}

func TestPath(t *testing.T) {
	keyfiles, err := filepath.Glob(path.Join(root, "var/db/xbps/keys/*.plist"))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range keyfiles {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var key PublicKey
		if err := ParsePublicKey(buf, &key); err != nil {
			t.Fatal(err)
		}
		res := key.Path(path.Join(root, "var/db/xbps"))
		if res != f {
			t.Errorf("Path() %q does not match %q", res, f)
		}
	}
}
