package uri

import (
	"testing"
)

var schemeTests = []struct {
	succeed bool
	rawuri  string
}{
	{true, "http://alpha.de.repo.voidlinux.org/"},
	{true, "https://alpha.de.repo.voidlinux.org/"},
	{true, "ftp://alpha.de.repo.voidlinux.org/"},
	{true, "/hostdir/binpkgs"},
	{true, "file:///hostdir/binpkgs"},
	{false, "torrent:///hostdir/binpkgs"},
}

func TestScheme(t *testing.T) {
	for _, tt := range schemeTests {
		_, err := Parse(tt.rawuri)
		if tt.succeed {
			if err != nil && tt.succeed {
				t.Fatalf("expected success for uri %q, got %s", tt.rawuri, err)
			}
		} else {
			if err == nil {
				t.Fatalf("expected error for uri %q", tt.rawuri)
			}
		}
	}
}

var isRemoteTests = []struct {
	remote bool
	rawuri string
}{
	{true, "http://alpha.de.repo.voidlinux.org/"},
	{true, "HTTP://alpha.de.repo.voidlinux.org/"},
	{true, "https://alpha.de.repo.voidlinux.org/"},
	{true, "ftp://alpha.de.repo.voidlinux.org/"},
	{false, "/hostdir/binpkgs"},
}

func TestIsRemote(t *testing.T) {
	for _, tt := range isRemoteTests {
		u, err := Parse(tt.rawuri)
		if err != nil {
			t.Fatal(err)
		}
		if u.IsRemote() != tt.remote {
			t.Fatalf("%#v.IsRemote returned %t expected %t", u, !tt.remote, tt.remote)
		}
	}
}

var cleanStringTests = []struct {
	rawuri string
	clean  string
}{
	{
		"http://alpha.de.repo.voidlinux.org/current/",
		"http___alpha_de_repo_voidlinux_org_current_",
	},
	{
		"HTTP://alpha.de.repo.voidlinux.org/current",
		"http___alpha_de_repo_voidlinux_org_current",
	},
	{
		"https://alpha.de.repo.voidlinux.org/current/",
		"https___alpha_de_repo_voidlinux_org_current_",
	},
	{
		"ftp://alpha.de.repo.voidlinux.org/current/",
		"ftp___alpha_de_repo_voidlinux_org_current_",
	},
}

func TestCleanString(t *testing.T) {
	for _, tt := range cleanStringTests {
		u, err := Parse(tt.rawuri)
		if err != nil {
			t.Fatal(err)
		}
		s := u.CacheString()
		if s != tt.clean {
			t.Errorf("%q returned %q expected %q", tt.rawuri, s, tt.clean)
		}
	}
}

func TestRepodata(t *testing.T) {
	res, err := Repodata("hostdir/binpkgs", "x86_64", "")
	if err != nil {
		t.Fatal(err)
	}
	expect := "hostdir/binpkgs/x86_64-repodata"
	if res != expect {
		t.Fatalf("expected %q, got %q\n", expect, res)
	}
	res, err = Repodata("https://repo-default.voidlinux.org/current", "x86_64", "/var/cache/xbps")
	if err != nil {
		t.Fatal(err)
	}
	expect = "/var/cache/xbps/https___repo-default_voidlinux_org_current/x86_64-repodata"
	if res != expect {
		t.Fatalf("expected %q, got %q\n", expect, res)
	}
}
