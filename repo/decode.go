package repo

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
	"howett.net/plist"
)

const (
	// IndexEntry is the file name of the index inside the repository data
	IndexEntry = "index.plist"
	// StageEntry is the file name of the stage inside the repository data
	StageEntry = "stage.plist"
	// MetaEntry is the file name of the metadata inside the repository data
	MetaEntry = "index-meta.plist"
)

// readCounter is a io.Reader wrapper that counts the number of bytes read
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

// Decoder is a repository data decoder
type Decoder struct {
	reader  readCounter
	decomp  *zstd.Decoder
	archive *tar.Reader
	header  *tar.Header
}

// Create a new repository data decoder
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

// ReadPlist reads the current repository entry and decodes its plist into v
func (dec *Decoder) ReadPlist(v any) error {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(dec.archive); err != nil {
		return err
	}
	rs := bytes.NewReader(buf.Bytes())
	plist := plist.NewDecoder(rs)
	if err := plist.Decode(v); err != nil {
		return err
	}
	return nil
}
