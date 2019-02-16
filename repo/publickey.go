package repo

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/DHowett/go-plist"
)

// plist structure of xbps generated key files.
type pubKey struct {
	Key      []byte `plist:"public-key"`
	Size     uint16 `plist:"public-key-size"`
	SignedBy string `plist:"signature-by"`
}

type PublicKey struct {
	Key      *rsa.PublicKey
	Size     uint16
	SignedBy string
}

func (p *PublicKey) UnmarshalPlist(unmarshal func(interface{}) error) error {
	var data pubKey
	if err := unmarshal(&data); err != nil {
		return err
	}
	block, _ := pem.Decode(data.Key)
	if block == nil {
		return errors.New("failed to decode PEM block")
	}
	if block.Type != "PUBLIC KEY" {
		return errors.New("not a public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		p.Key = pub
		return nil
	default:
		return errors.New("unsupported public key type")
	}
	return nil
}

// TODO:
func (p *PublicKey) Filename() string {
	return ""
}

func ParsePublicKey(data []byte, key *PublicKey) error {
	if _, err := plist.Unmarshal(data, key); err != nil {
		return err
	}
	return nil
}
