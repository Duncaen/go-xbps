package util

import (
	"testing"
)

func TestFileSha256(t *testing.T) {
	hash, err := FileSha256("/etc/os-release")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}
