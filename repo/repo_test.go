package repo

import (
	"testing"
)

func TestUnstaged(t *testing.T) {
	var r *Repository
	var err error
	if r, err = Open("/tmp/test-repo2", "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Packages {
		t.Logf("%v: %v", k, v)
	}
	for k, v := range r.StagedPackages {
		t.Logf("%v: %v", k, v)
	}
	if r.Staged {
		t.Fatal("repo is staged")
	}
}

func TestStaged(t *testing.T) {
	var r *Repository
	var err error
	if r, err = Open("/tmp/test-repo", "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Packages {
		t.Logf("%v: %v", k, v)
	}
	for k, v := range r.StagedPackages {
		t.Logf("%v: %v", k, v)
	}
	if ! r.Staged {
		t.Fatal("repo is not staged")
	}
}
