package envcrypto

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/nacl/box"
)

type Encoding interface {
	EncodeToString([]byte) string
	DecodeString(string) ([]byte, error)
}

type Box struct {
	source Source

	// options
	rand                     io.Reader
	encoding                 Encoding
	publicKeyKey             string
	privateKeyKey            string
	encryptedDataValuePrefix string

	// state
	publicKeyString  string
	privateKeyString string
	publicKey        *[32]byte
	privateKey       *[32]byte
}

const sopsEncryptedPrefix = "ENC["

func newBox(source Source, options ...Option) *Box {
	b := &Box{
		source: source,

		rand:     rand.Reader,
		encoding: base32.StdEncoding.WithPadding(base32.NoPadding),

		publicKeyKey:  "ENVCRYPTO_PUBLIC_KEY",
		privateKeyKey: "ENVCRYPTO_PRIVATE_KEY",

		encryptedDataValuePrefix: "!envcrypto:",
	}
	for _, option := range options {
		option.applyToBox(b)
	}
	return b
}

func New(options ...Option) (MapSource, error) {
	s := make(MapSource)
	b := newBox(s, options...)
	publicKey, privateKey, err := box.GenerateKey(b.rand)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	b.publicKey, b.privateKey = publicKey, privateKey
	s[b.publicKeyKey] = b.encoding.EncodeToString(b.publicKey[:])
	s[b.privateKeyKey] = b.encoding.EncodeToString(b.privateKey[:])
	return s, nil
}

func (b *Box) parseKey(s string) (*[32]byte, error) {
	key, err := b.encoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("key is not 32 bytes long")
	}
	var key32 [32]byte
	copy(key32[:], key)
	return &key32, nil
}

func Open(source Source, options ...Option) (*Box, error) {
	b := newBox(source, options...)

	if b.publicKeyString == "" {
		b.publicKeyString, _ = source.Lookup(b.publicKeyKey)
	}
	if b.publicKeyString == "" {
		return nil, fmt.Errorf("no public key")
	}
	publicKey, err := b.parseKey(b.publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %q: %w", b.publicKeyString, err)
	}
	b.publicKey = publicKey

	if b.privateKeyString == "" {
		b.privateKeyString, _ = source.Lookup(b.privateKeyKey)
	}
	if b.privateKeyString != "" && !strings.HasPrefix(b.privateKeyString, sopsEncryptedPrefix) {
		privateKey, err := b.parseKey(b.privateKeyString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		b.privateKey = privateKey
	}

	return b, nil
}

func (b *Box) Encrypt(value string) (string, error) {
	out, err := box.SealAnonymous(nil, []byte(value), b.publicKey, b.rand)
	if err != nil {
		return "", err
	}
	encryptedValue := b.encoding.EncodeToString(out)
	return b.encryptedDataValuePrefix + encryptedValue, nil
}

func (b *Box) Decrypt(encryptedValue string) (string, error) {
	encryptedValue = strings.TrimPrefix(encryptedValue, b.encryptedDataValuePrefix)
	encryptedValueBytes, err := b.encoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}
	if b.publicKey == nil {
		return "", fmt.Errorf("no public key")
	}
	if b.privateKey == nil {
		return "", fmt.Errorf("no private key")
	}
	out, ok := box.OpenAnonymous(nil, encryptedValueBytes, b.publicKey, b.privateKey)
	if !ok {
		return "", fmt.Errorf("failed to decrypt value")
	}
	return string(out), nil
}

func (b *Box) Get(key string) (string, error) {
	if value, ok := b.source.Lookup(key); ok {
		if b.encryptedDataValuePrefix != "" && strings.HasPrefix(value, b.encryptedDataValuePrefix) {
			return b.Decrypt(value)
		}
		return value, nil
	}
	return "", fmt.Errorf("no value for key %q", key)
}

func (b *Box) All() (map[string]string, error) {
	m := make(map[string]string)
	for _, key := range b.source.Keys() {
		if key == b.publicKeyKey || key == b.privateKeyKey {
			continue
		}
		value, err := b.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %q: %w", key, err)
		}
		m[key] = value
	}
	return m, nil
}
