package uri

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type URI url.URL

func Parse(rawuri string) (*URI, error) {
	u, err := url.Parse(rawuri)
	if err != nil {
		return nil, err
	}
	uri := (*URI)(u)
	if _, err := uri.isSupported(); err != nil {
		return nil, err
	}
	return uri, nil
}

// supported remote repository schemes
var remoteSchemes = []string{
	"ftp",
	"http",
	"https",
}

// supported local repository schemes
var localSchemes = []string{
	"file",
	"",
}

func (u *URI) isSupported() (bool, error) {
	if slices.Contains(remoteSchemes, u.Scheme) {
		return true, nil
	}
	if slices.Contains(localSchemes, u.Scheme) {
		return true, nil
	}
	return false, fmt.Errorf("scheme is not supported: %q", u.Scheme)
}

// IsRemote returns true if the repository URI is remote
func (u *URI) IsRemote() bool {
	return slices.Contains(remoteSchemes, u.Scheme)
}

// String returns the the url as string
func (u *URI) String() string {
	return (*url.URL)(u).String()
}

// Directory returns
func (u *URI) CleanString() string {
	return strings.NewReplacer(".", "_", "/", "_", ":", "_").Replace(u.String())
}
