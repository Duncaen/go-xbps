package repo

import (
	"testing"
)

func TestUnstaged(t *testing.T) {
	var r *Repository
	var err error
	if r, err = Open("/home/duncan/repos/go-xbps/testrepo", "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Packages {
		t.Logf("%v: %v", k, v)
	}
}

func TestStaged(t *testing.T) {
	var r *Repository
	var err error
	if r, err = Open("/home/duncan/repos/go-xbps/testrepo-staged", "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Packages {
		t.Logf("%v: %v", k, v)
	}
}
