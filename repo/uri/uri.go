package uri

import (
	"fmt"
	"net/url"
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
	"",
	"file",
}

func (u *URI) isSupported() (bool, error) {
	ok := false
	for _, scheme := range append(remoteSchemes, localSchemes...) {
		if u.Scheme == scheme {
			ok = true
			break
		}
	}
	if !ok {
		return ok, fmt.Errorf("scheme %q is not supported", u.Scheme)
	}
	return ok, nil
}

// IsRemote returns true if the repository URI is remote
func (u *URI) IsRemote() bool {
	for _, scheme := range remoteSchemes {
		if u.Scheme == scheme {
			return true
		}
	}
	return false
}

func (u *URI) String() string {
	return (*url.URL)(u).String()
}

// Directory returns
func (u *URI) CleanString() string {
	return strings.NewReplacer(".", "_", "/", "_", ":", "_").Replace(u.String())
}
