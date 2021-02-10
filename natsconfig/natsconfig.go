// Package natsconfig provides config and flags for connecting to NATS.
package natsconfig

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/spf13/pflag"
)

// AuthConfig represents NATS authentication configuration.
type AuthConfig struct {
	Username        string `json:"username,omitempty" yaml:"username,omitempty"`
	Password        string `json:"password,omitempty" yaml:"password,omitempty"`
	PasswordFile    string `json:"passwordFile,omitempty" yaml:"passwordFile,omitempty"`
	CredentialsFile string `json:"credentialsFile,omitempty" yaml:"credentialsFile,omitempty"`
	JWTFile         string `json:"jwtFile,omitempty" yaml:"jwtFile,omitempty"`
	SeedFile        string `json:"seedFile,omitempty" yaml:"seedFile,omitempty"`
}

// DefaultAuthConfig returns the default configuration for NATS authentication.
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{}
}

// Flags returns a flagset that can be added to the command line.
func (c *AuthConfig) Flags(prefix string, defaults *AuthConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultAuthConfig()
	}
	flags.StringVar(&c.Username, prefix+"auth.username", defaults.Username, "NATS username")
	flags.StringVar(&c.Password, prefix+"auth.password", defaults.Password, "NATS password")
	flags.StringVar(&c.PasswordFile, prefix+"auth.passwordFile", defaults.PasswordFile, "NATS password file")
	flags.StringVar(&c.CredentialsFile, prefix+"auth.credentialsFile", defaults.CredentialsFile, "NATS credentials file")
	flags.StringVar(&c.JWTFile, prefix+"auth.jwtFile", defaults.JWTFile, "NATS JWT file")
	flags.StringVar(&c.SeedFile, prefix+"auth.seedFile", defaults.SeedFile, "NATS seed file")
	return &flags
}

// Load loads the NATS credentials.
func (c *AuthConfig) Load(_ context.Context) (nats.Option, error) {
	if c.Username != "" {
		if c.PasswordFile != "" {
			passwordBytes, err := ioutil.ReadFile(c.PasswordFile)
			if err != nil {
				return nil, err
			}
			c.Password = string(bytes.TrimSpace(passwordBytes))
		}
		return func(opts *nats.Options) error {
			opts.User = c.Username
			opts.Password = c.Password
			return nil
		}, nil
	}
	if c.CredentialsFile != "" {
		credentials, err := ioutil.ReadFile(c.CredentialsFile)
		if err != nil {
			return nil, err
		}
		jwt, err := nkeys.ParseDecoratedJWT(credentials)
		if err != nil {
			return nil, err
		}
		nkey, err := nkeys.ParseDecoratedNKey(credentials)
		if err != nil {
			return nil, err
		}
		return nats.UserJWT(
			func() (string, error) {
				return jwt, nil
			},
			func(input []byte) ([]byte, error) {
				return nkey.Sign(input)
			},
		), nil
	}
	if c.JWTFile != "" && c.SeedFile != "" {
		jwtBytes, err := ioutil.ReadFile(c.JWTFile)
		if err != nil {
			return nil, err
		}
		jwt := string(jwtBytes)
		seedBytes, err := ioutil.ReadFile(c.SeedFile)
		if err != nil {
			return nil, err
		}
		nkey, err := nkeys.FromSeed(seedBytes)
		if err != nil {
			return nil, err
		}
		return nats.UserJWT(
			func() (string, error) {
				return jwt, nil
			},
			func(input []byte) ([]byte, error) {
				return nkey.Sign(input)
			},
		), nil
	}
	return func(*nats.Options) error { return nil }, nil
}

// Config is the configuration for the NATS connection.
type Config struct {
	Servers   []string    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Name      string      `json:"name,omitempty" yaml:"name,omitempty"`
	Auth      AuthConfig  `json:"auth,omitempty" yaml:"auth,omitempty"`
	TLSConfig *tls.Config `json:"-" yaml:"-"`
}

// DefaultConfig returns the default configuration for the NATS connection.
func DefaultConfig() *Config {
	return &Config{
		Servers: []string{nats.DefaultURL},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string, defaults *Config) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultConfig()
	}
	flags.StringSliceVar(&c.Servers, prefix+"servers", defaults.Servers, "NATS servers")
	flags.StringVar(&c.Name, prefix+"name", defaults.Name, "Name to send to the NATS servers")
	flags.AddFlagSet(c.Auth.Flags(prefix, &defaults.Auth))
	return &flags
}

// Connect connects to NATS using this configuration.
func (c *Config) Connect(ctx context.Context) (*nats.Conn, error) {
	opts := nats.GetDefaultOptions()
	opts.Servers = c.Servers
	opts.Name = c.Name
	if c.TLSConfig != nil {
		opts.Secure = true
		opts.TLSConfig = c.TLSConfig
	}
	authOption, err := c.Auth.Load(ctx)
	if err != nil {
		return nil, err
	}
	if err = authOption(&opts); err != nil {
		return nil, err
	}
	return opts.Connect()
}
