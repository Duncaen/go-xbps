package util

import (
	"crypto/sha256"
	"io"
	"os"
)

func FileSha256(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
