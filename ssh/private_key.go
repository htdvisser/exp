package ssh

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// PrivateKeyConfig is the configuration of an SSH private key.
type PrivateKeyConfig struct {
	File       string `json:"file" yaml:"file"`
	Passphrase string `json:"passphrase" yaml:"passphrase"`
}

// Validate validates the configuration and returns an error if it is not valid.
func (c PrivateKeyConfig) Validate() error {
	if c.File == "" {
		return fmt.Errorf("missing private key file in PrivateKeyConfig")
	}
	return nil
}

// Build builds an ssh.Signer from the configuration.
func (c PrivateKeyConfig) Build() (ssh.Signer, error) {
	pemBytes, err := ioutil.ReadFile(c.File)
	if err != nil {
		return nil, fmt.Errorf("private key %q could not be read: %w", c.File, err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("private key %q could not be parsed: %w", c.File, err)
	}
	return signer, nil
}
