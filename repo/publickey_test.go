package repo

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
)

var root = "/"

func TestFile(t *testing.T) {
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
	}
}
