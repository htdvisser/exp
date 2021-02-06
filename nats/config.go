package nats

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/spf13/pflag"
)

// FileReader reads a file and returns its bytes.
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// CredsConfig represents NATS credentials.
type CredsConfig struct {
	CredsFile string
	JWTFile   string
	SeedFile  string
}

// Load loads the NATS credentials. If fileReader is nil, this uses ioutil.ReadFile.
func (f CredsConfig) Load(_ context.Context, fileReader FileReader) (nats.Option, error) {
	readFile := ioutil.ReadFile
	if fileReader != nil {
		readFile = fileReader.ReadFile
	}
	if f.CredsFile != "" {
		log.Printf("Load NATS credentials from %q", f.CredsFile)
		credentials, err := readFile(f.CredsFile)
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
	if f.JWTFile != "" && f.SeedFile != "" {
		log.Printf("Load NATS credentials from %q and %q", f.JWTFile, f.SeedFile)
		jwtBytes, err := readFile(f.JWTFile)
		if err != nil {
			return nil, err
		}
		jwt := string(jwtBytes)
		seedBytes, err := readFile(f.SeedFile)
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
	Servers      []string
	Name         string
	Username     string
	Password     string
	PasswordFile string
	Store        string
	Credentials  CredsConfig
	TLSConfig    *tls.Config
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
	flags.StringVar(&c.Username, prefix+"auth.username", defaults.Username, "NATS username")
	flags.StringVar(&c.Password, prefix+"auth.password", defaults.Password, "NATS password")
	flags.StringVar(&c.PasswordFile, prefix+"auth.password-file", defaults.PasswordFile, "NATS password file")
	flags.StringVar(&c.Store, prefix+"store", defaults.Store, "NATS credentials store")
	flags.StringVar(&c.Credentials.CredsFile, prefix+"auth.credentials-file", defaults.Credentials.CredsFile, "NATS credentials file")
	flags.StringVar(&c.Credentials.JWTFile, prefix+"auth.jwt-file", defaults.Credentials.JWTFile, "NATS JWT file")
	flags.StringVar(&c.Credentials.SeedFile, prefix+"auth.seed-file", defaults.Credentials.SeedFile, "NATS seed file")
	return &flags
}

func (c *Config) getStore() (FileReader, error) {
	switch c.Store {
	case "", "file":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported TLS certificate store %q", c.Store)
	}
}

// Connect opens a connection to NATS.
func (c *Config) Connect(ctx context.Context) (*nats.Conn, error) {
	opts := nats.GetDefaultOptions()
	opts.Servers = c.Servers
	opts.Name = c.Name
	opts.User = c.Username
	opts.Password = c.Password
	if c.PasswordFile != "" {
		passwordBytes, err := ioutil.ReadFile(c.PasswordFile)
		if err != nil {
			return nil, err
		}
		opts.Password = string(bytes.TrimSpace(passwordBytes))
	}
	if c.TLSConfig != nil {
		opts.Secure = true
		opts.TLSConfig = c.TLSConfig
	}
	store, err := c.getStore()
	if err != nil {
		return nil, err
	}
	opt, err := c.Credentials.Load(ctx, store)
	if err != nil {
		return nil, err
	}
	if err = opt(&opts); err != nil {
		return nil, err
	}
	log.Printf("Connect to NATS %v", opts.Servers)
	return opts.Connect()
}
