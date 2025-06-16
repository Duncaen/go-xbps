package repo

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
	"howett.net/plist"
)

// io.Reader wrapper that counts the number of bytes read
type readCounter struct {
	io.Reader
	n int64
}

// Read implementation that counts bytes read
func (counter *readCounter) Read(p []byte) (int, error) {
	n, err := counter.Reader.Read(p)
	counter.n += int64(n)
	return n, err
}

// Decoder is a repository data reader
type Decoder struct {
	reader  readCounter
	decomp  *zstd.Decoder
	archive *tar.Reader
	header  *tar.Header
}

// Create a new repository data reader
func NewDecoder(r io.Reader) (*Decoder, error) {
	var err error
	dec := &Decoder{
		reader: readCounter{r, 0},
	}
	dec.decomp, err = zstd.NewReader(&dec.reader)
	if err != nil {
		return nil, err
	}
	dec.archive = tar.NewReader(dec.decomp)
	return dec, nil
}

// Close the repository data reader
func (dec *Decoder) Close() {
	dec.decomp.Close()
}

// Next returns the name of the next file in the repository data
func (dec *Decoder) Next() (string, error) {
	var err error
	dec.header, err = dec.archive.Next()
	if err != nil {
		return "", err
	}
	return dec.header.Name, nil
}

// ReadPackages decodes the current package dictionary inside the repository data
func (dec *Decoder) ReadPackages() (map[string]Package, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(dec.archive); err != nil {
		return nil, err
	}
	rs := bytes.NewReader(buf.Bytes())
	plist := plist.NewDecoder(rs)
	pkgs := make(map[string]Package)
	if err := plist.Decode(&pkgs); err != nil {
		return nil, err
	}
	return pkgs, nil
}
