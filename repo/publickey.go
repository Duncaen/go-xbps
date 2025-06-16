package repo

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"howett.net/plist"
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
	p.Size, p.SignedBy = data.Size, data.SignedBy
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

// Returns the path where xbps would store the public key
func (p *PublicKey) Path(dbdir string) string {
	return filepath.Join(dbdir, "keys", p.Filename())
}

// Returns the filename xbps would use to store the public key
func (p *PublicKey) Filename() string {
	return fmt.Sprintf("%s.plist", p.Fingerprint())
}

// Returns an OpenSSH compatible public key fingerprint
func (p *PublicKey) Fingerprint() string {
	// BUG(duncaen): xbps uses md5 fingerprints, this is the old format openssh used
	pubKey, err := ssh.NewPublicKey(p.Key)
	if err != nil {
		// this should never happen with rsa keys
		panic(err)
	}
	return ssh.FingerprintLegacyMD5(pubKey)
}

func ParsePublicKey(data []byte, key *PublicKey) error {
	if _, err := plist.Unmarshal(data, key); err != nil {
		return err
	}
	return nil
}
