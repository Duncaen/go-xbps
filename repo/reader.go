package repo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"

	"howett.net/plist"
)

type Reader struct {
	gzr *gzip.Reader
	tr *tar.Reader
}

func NewReader(r io.Reader) (*Reader, error) {
	var err error
	rd := &Reader{}
	rd.gzr, err = gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	rd.tr = tar.NewReader(rd.gzr)
	return rd, nil
}

func (rd *Reader) Close() {
	rd.gzr.Close()
}

func (rd *Reader) Next() (string, error) {
	hdr, err := rd.tr.Next()
	if err != nil {
		return "", err
	}
	return hdr.Name, nil
}

func (rd *Reader) ReadPackages() (map[string]Package, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(rd.tr); err != nil {
		return nil, err
	}
	rs := bytes.NewReader(buf.Bytes())
	dec := plist.NewDecoder(rs)
	pkgs := make(map[string]Package)
	if err := dec.Decode(&pkgs); err != nil {
		return nil, err
	}
	return pkgs, nil
}

func (rd *Reader) ReadPublicKey() PublicKey {
	return PublicKey{}
}
