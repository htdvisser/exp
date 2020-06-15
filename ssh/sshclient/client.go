// Package sshclient provides configuration structs for constructing SSH clients.
package sshclient

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	hssh "htdvisser.dev/exp/ssh"
	"htdvisser.dev/exp/ssh/aws"
)

func stringFallback(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func durationFallback(d, fallback time.Duration) time.Duration {
	if d != 0 {
		return d
	}
	return fallback
}

// HostKeyConfig is the configuration for host key verification.
type HostKeyConfig struct {
	Source     string `json:"source" yaml:"source"`
	KnownHosts struct {
		File string `json:"file" yaml:"file"`
	} `json:"known_hosts" yaml:"known_hosts"`
}

// Validate validates the configuration and returns an error if it is not valid.
func (c HostKeyConfig) Validate() error {
	switch c.Source {
	case "insecure_ignore":
	case "known_hosts":
		if c.KnownHosts.File == "" {
			return fmt.Errorf("missing known_hosts file in HostKeyConfig")
		}
	default:
		return fmt.Errorf("invalid source %q for HostKeyConfig", c.Source)
	}
	return nil
}

func (c HostKeyConfig) build() (ssh.HostKeyCallback, error) {
	switch c.Source {
	case "insecure_ignore":
		return ssh.InsecureIgnoreHostKey(), nil
	case "known_hosts":
		cb, err := knownhosts.New(c.KnownHosts.File)
		if err != nil {
			return nil, fmt.Errorf("known_hosts file %q failed to load: %w", c.KnownHosts.File, err)
		}
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return cb(hostname, remote, key)
		}, nil
	}
	return nil, nil
}

// BannerConfig is the configuration for handling the banner received from the SSH server.
type BannerConfig struct {
	// TODO: Define handler types.
}

// Validate validates the configuration and returns an error if it is not valid.
func (BannerConfig) Validate() error {
	// TODO: Implement validator.
	return nil
}

func (BannerConfig) build() (ssh.BannerCallback, error) {
	// TODO: Implement callbacks.
	return nil, nil
}

// AuthMethodConfig is the configuration for the authentication method.
type AuthMethodConfig struct {
	Method      string                  `json:"method" yaml:"method"`
	Password    string                  `json:"password" yaml:"password"`
	PrivateKeys []hssh.PrivateKeyConfig `json:"private_keys" yaml:"private_keys"`
	AWSKMS      aws.KMSConfig           `json:"aws_kms" yaml:"aws_kms"`
	// TODO: Support GCP Cloud HSM.
}

// Validate validates the configuration and returns an error if it is not valid.
func (c AuthMethodConfig) Validate() error {
	switch c.Method {
	case "password":
		if c.Password == "" {
			return fmt.Errorf("missing password file in AuthMethodConfig")
		}
	case "private_keys":
		if len(c.PrivateKeys) == 0 {
			return fmt.Errorf("missing private keys in AuthMethodConfig")
		}
		for _, pk := range c.PrivateKeys {
			if err := pk.Validate(); err != nil {
				return fmt.Errorf("invalid private key in AuthMethodConfig: %w", err)
			}
		}
	case "aws_kms":
		if err := c.AWSKMS.Validate(); err != nil {
			return fmt.Errorf("invalid AWS KMS in AuthMethodConfig: %w", err)
		}
	default:
		return fmt.Errorf("invalid method %q for AuthMethodConfig", c.Method)
	}
	return nil
}

func (c AuthMethodConfig) build() (ssh.AuthMethod, error) {
	switch c.Method {
	case "password":
		return ssh.PasswordCallback(func() (string, error) {
			return c.Password, nil
		}), nil
	case "private_keys":
		return ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
			var (
				signers = make([]ssh.Signer, len(c.PrivateKeys))
				err     error
			)
			for i, pkc := range c.PrivateKeys {
				signers[i], err = pkc.Build()
				if err != nil {
					return nil, fmt.Errorf("failed to build signer for private key %d: %w", i, err)
				}
			}
			return signers, nil
		}), nil
	case "aws_kms":
		return ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
			signer, err := c.AWSKMS.Build()
			if err != nil {
				return nil, fmt.Errorf("failed to build AWS KMS signer: %w", err)
			}
			return []ssh.Signer{signer}, nil
		}), nil
	}
	return nil, nil
}

// ConnectConfig is the configuration for connecting to an SSH server.
type ConnectConfig struct {
	Network       string             `json:"network,omitempty" yaml:"network,omitempty"`
	Address       string             `json:"address" yaml:"address"`
	Timeout       time.Duration      `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	KeepAlive     time.Duration      `json:"keep_alive,omitempty" yaml:"keep_alive,omitempty"`
	HostKey       HostKeyConfig      `json:"host_key,omitempty" yaml:"host_key,omitempty"`
	Banner        BannerConfig       `json:"banner,omitempty" yaml:"banner,omitempty"`
	ClientVersion string             `json:"client_version,omitempty" yaml:"client_version,omitempty"`
	Username      string             `json:"username" yaml:"username"`
	AuthMethods   []AuthMethodConfig `json:"auth_methods" yaml:"auth_methods"`
}

// Validate validates the configuration and returns an error if it is not valid.
func (c ConnectConfig) Validate() error {
	switch c.Network {
	case "", "tcp", "tcp4", "tcp6":
		if _, err := net.ResolveTCPAddr(c.Network, c.Address); err != nil {
			return fmt.Errorf("invalid address %q: %w", c.Address, err)
		}
	case "unix":
		info, err := os.Stat(c.Address)
		if err != nil {
			return fmt.Errorf("invalid address %q: %w", c.Address, err)
		}
		if info.Mode()&os.ModeSocket != os.ModeSocket {
			return fmt.Errorf("address %q does not seem to be a socket", c.Address)
		}
	default:
		return fmt.Errorf("invalid network %q for ConnectConfig", c.Network)
	}
	if err := c.HostKey.Validate(); err != nil {
		return fmt.Errorf("invalid host key in ConnectConfig: %w", err)
	}
	if err := c.Banner.Validate(); err != nil {
		return fmt.Errorf("invalid banner in ConnectConfig: %w", err)
	}
	if c.Username == "" {
		return fmt.Errorf("missing username in ConnectConfig")
	}
	if len(c.AuthMethods) == 0 {
		return fmt.Errorf("missing auth methods in ConnectConfig")
	}
	for _, am := range c.AuthMethods {
		if err := am.Validate(); err != nil {
			return fmt.Errorf("invalid auth method in ConnectConfig: %w", err)
		}
	}
	return nil
}

// Dial dials the configured SSH server.
func (c ConnectConfig) Dial(ctx context.Context) (*ssh.Client, error) {
	var (
		authMethods = make([]ssh.AuthMethod, len(c.AuthMethods))
		err         error
	)
	for i, amc := range c.AuthMethods {
		if authMethods[i], err = amc.build(); err != nil {
			return nil, fmt.Errorf("failed to build auth method %d: %w", i, err)
		}
	}
	hostKeyCallback, err := c.HostKey.build()
	if err != nil {
		return nil, fmt.Errorf("failed to build host key callback: %w", err)
	}
	bannerCallback, err := c.Banner.build()
	if err != nil {
		return nil, fmt.Errorf("failed to build banner callback: %w", err)
	}
	d := net.Dialer{
		Timeout:   durationFallback(c.Timeout, 10*time.Second),
		KeepAlive: durationFallback(c.KeepAlive, 10*time.Second),
	}
	tcpConn, err := d.DialContext(
		ctx,
		stringFallback(c.Network, "tcp"),
		c.Address,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %q: %w", c.Address, err)
	}
	sshConn, sshChannels, sshRequests, err := ssh.NewClientConn(tcpConn, c.Address, &ssh.ClientConfig{
		Config:          ssh.Config{},
		Timeout:         durationFallback(c.Timeout, 10*time.Second),
		HostKeyCallback: hostKeyCallback,
		BannerCallback:  bannerCallback,
		ClientVersion:   c.ClientVersion,
		User:            c.Username,
		Auth:            authMethods,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH client connection: %w", err)
	}
	return ssh.NewClient(sshConn, sshChannels, sshRequests), nil
}
