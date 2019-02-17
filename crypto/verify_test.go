package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Duncaen/go-xbps/repo"
	"github.com/Duncaen/go-xbps/util"
)

var root = "/"

func readAllKeys() ([]repo.PublicKey, error) {
	keyfiles, err := filepath.Glob(path.Join(root, "var/db/xbps/keys/*.plist"))
	if err != nil {
		return nil, err
	}
	keys := make([]repo.PublicKey, len(keyfiles))
	for i, f := range keyfiles {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if err := repo.ParsePublicKey(buf, &keys[i]); err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func TestVerify(t *testing.T) {
	keys, err := readAllKeys()
	if err != nil {
		t.Skipf("skipping test because of: %s", err)
	}
	sigs, err := filepath.Glob(path.Join(root, "var/cache/xbps/*.xbps.sig"))
	if err != nil {
		t.Skipf("skipping test because of: %s", err)
	}
	i := 0
	for _, sigf := range sigs {
		if i > 10 {
			break
		}
		hash, err := util.FileSha256(strings.TrimSuffix(sigf, ".sig"))
		if err != nil {
			continue
		}
		sig, err := ioutil.ReadFile(sigf)
		if err != nil {
			t.Log(err)
			continue
		}
		i++
		for _, k := range keys {
			if err := Verify(k.Key, hash, sig); err != nil {
				t.Logf("key %s: %q: %s", k.Fingerprint(), sigf, err)
				continue
			}
			t.Logf("key %s: %q: sucessfully verified", k.Fingerprint(), sigf)
		}
	}
}

func TestSign(t *testing.T) {
	msg := "Hello World"
	hashed := sha256.Sum256([]byte(msg))
	priv, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := Sign(priv, hashed[:])
	if err != nil {
		t.Fatal(err)
	}
	pub := priv.Public().(*rsa.PublicKey)
	if err := Verify(pub, hashed[:], sig); err != nil {
		t.Fatal(err)
	}
}
