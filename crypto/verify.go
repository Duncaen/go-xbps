// Package crypto implements XBPS' signature signing and verification using crypto/rsa.
//
// XBPS SHA1-SHA256 workaround
//
// xbps specifies sha1 as algorithm for RSA_{sign,verify} but uses
// a sha256 message and length.
//
// xbps still uses full sha256 hashes for the verification, the describes issue
// is just a parsing issue which with strict ASN1/PKCS1 implementations.
//
// As a workaround this implementation disables golangs PKCS1v15 prefix by using
// crypto.Hash(0) as hash argument for rsa.VerifyPKCS1v15 and rsa.SignPKCS1v15
// and uses a precomputed ASN1 prefix instead.
//
// The ANS1 prefix comes from xbps generated signature, dumped using the
// openssl command:
//  openssl rsautl -verify -in foo-1.0_1.noarch.xbps.sig  -inkey ~/.ssh/id_rsa -raw -hexdump
//
// XBPS implementation:
//
// https://github.com/void-linux/xbps/blob/b4eebaf/lib/verifysig.c#L66
//
// Original bug report:
//
// https://github.com/voidlinux/xbps/issues/146
//
// Note: golang also hardcodes the ASN1 prefix for performance reasons:
//
// https://github.com/golang/go/blob/dca707b/src/crypto/rsa/pkcs1v15.go#L210
package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

var errHashMismatch = errors.New("input must be a sha256 hash")

// ASN1 message prefix with SHA1 algorithm id and sha256 digest length
var sha1x256Prefix = []byte{0x30, 0x2d, 0x30, 0x09, 0x06, 0x05, 0x2b, 0x0e, 0x03, 0x02, 0x1a, 0x05, 0x00, 0x04, 0x20}

// Verify a sha256 hash signature and work around xbps' quirks
func Verify(pub *rsa.PublicKey, hashed []byte, sig []byte) error {
	tLen := len(hashed)
	if tLen != 32 {
		return errHashMismatch
	}
	pLen := len(sha1x256Prefix)
	t := make([]byte, pLen+tLen)
	copy(t[:pLen], sha1x256Prefix)
	copy(t[pLen:], hashed)
	return rsa.VerifyPKCS1v15(pub, crypto.Hash(0), t, sig)
}

// Sign a sha256 hash using xbps' quirks
func Sign(priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	if tLen != 32 {
		return nil, errHashMismatch
	}
	pLen := len(sha1x256Prefix)
	t := make([]byte, pLen+tLen)
	copy(t[:pLen], sha1x256Prefix)
	copy(t[pLen:], hashed)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.Hash(0), t)
}
