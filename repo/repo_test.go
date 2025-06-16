package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

func createPackage(pkgver string, repodir string, destdir string, extra_args ...string) error {
	if err := os.MkdirAll(repodir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(destdir, 0o755); err != nil {
		return err
	}
	args := []string{
		"-A", "noarch",
		"-n", pkgver,
		"-s", "test package",
	}
	args = append(args, extra_args...)
	args = append(args, destdir)
	cmd := exec.Command("xbps-create", args...)
	cmd.Dir = repodir
	cmd.Env = append(cmd.Env, "XBPS_ARCH=x86_64")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("xbps-rindex", "-a", path.Join(repodir, fmt.Sprintf("%s.noarch.xbps", pkgver)))
	cmd.Env = append(cmd.Env, "XBPS_ARCH=x86_64")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func TestUnstaged(t *testing.T) {
	dir := t.TempDir()
	pkgdir := path.Join(dir, "pkg")
	repodir := path.Join(dir, "repo")
	if err := createPackage("foo-1.0_1", repodir, pkgdir, "--shlib-provides", "libfoo.so.1"); err != nil {
		t.Fatal("failed to create package:", err)
	}
	if err := createPackage("bar-1.0_1", repodir, pkgdir, "--shlib-requires", "libfoo.so.1"); err != nil {
		t.Fatal("failed to create package:", err)
	}

	var r *Repository
	var err error
	if r, err = Open(repodir, "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Index {
		t.Logf("%v: %v", k, v)
	}
	for k, v := range r.Stage {
		t.Logf("%v: %v", k, v)
	}
	if len(r.Stage) > 0 {
		t.Fatal("repo is staged")
	}
}

func TestStaged(t *testing.T) {
	dir := t.TempDir()
	pkgdir := path.Join(dir, "pkg")
	repodir := path.Join(dir, "repo")
	if err := createPackage("foo-1.0_1", repodir, pkgdir, "--shlib-provides", "libfoo.so.1"); err != nil {
		t.Fatal("failed to create package:", err)
	}
	if err := createPackage("bar-1.0_1", repodir, pkgdir, "--shlib-requires", "libfoo.so.1"); err != nil {
		t.Fatal("failed to create package:", err)
	}
	if err := createPackage("foo-2.0_1", repodir, pkgdir, "--shlib-provides", "libfoo.so.2"); err != nil {
		t.Fatal("failed to create package:", err)
	}

	var r *Repository
	var err error
	if r, err = Open(repodir, "x86_64"); err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	for k, v := range r.Index {
		t.Logf("%v: %v", k, v)
	}
	for k, v := range r.Stage {
		t.Logf("%v: %v", k, v)
	}
	if len(r.Stage) == 0 {
		t.Fatal("repo is not staged")
	}
}
