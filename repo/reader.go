package repo

import (
	"archive/tar"
	"bytes"
	"io"

	"howett.net/plist"
	"github.com/klauspost/compress/zstd"
)

type Reader struct {
	crd *zstd.Decoder
	tr *tar.Reader
}

func NewReader(r io.Reader) (*Reader, error) {
	var err error
	rd := &Reader{}
	rd.crd, err = zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	rd.tr = tar.NewReader(rd.crd)
	return rd, nil
}

func (rd *Reader) Close() {
	rd.crd.Close()
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
