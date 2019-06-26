package repo

import (
	"testing"
)

func TestFromFile(t *testing.T) {
	r := &Repository{}
	if err := r.Open("/var/db/xbps/https___alpha_de_repo_voidlinux_org_current/x86_64-repodata"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Packages {
		t.Logf("%v: %v", k, v)
	}
}
