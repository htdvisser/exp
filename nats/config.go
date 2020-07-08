package nats

import (
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
	Creds string
	JWT   string
	Seed  string
}

// Load loads the NATS credentials. If fileReader is nil, this uses ioutil.ReadFile.
func (f CredsConfig) Load(_ context.Context, fileReader FileReader) (nats.Option, error) {
	readFile := ioutil.ReadFile
	if fileReader != nil {
		readFile = fileReader.ReadFile
	}
	if f.Creds != "" {
		log.Printf("Load NATS credentials from %q", f.Creds)
		credentials, err := readFile(f.Creds)
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
	if f.JWT != "" && f.Seed != "" {
		log.Printf("Load NATS credentials from %q and %q", f.JWT, f.Seed)
		jwtBytes, err := readFile(f.JWT)
		if err != nil {
			return nil, err
		}
		jwt := string(jwtBytes)
		seedBytes, err := readFile(f.Seed)
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
	Servers     []string
	Name        string
	Username    string
	Password    string
	Store       string
	Credentials CredsConfig
	TLSConfig   *tls.Config
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string) *pflag.FlagSet {
	var flags pflag.FlagSet
	flags.StringSliceVar(&c.Servers, prefix+"servers", []string{nats.DefaultURL}, "NATS servers")
	flags.StringVar(&c.Name, prefix+"name", "", "Name to send to the NATS servers")
	flags.StringVar(&c.Username, prefix+"auth.username", "", "NATS username")
	flags.StringVar(&c.Password, prefix+"auth.password", "", "NATS password")
	flags.StringVar(&c.Store, prefix+"store", "", "NATS credentials store")
	flags.StringVar(&c.Credentials.Creds, prefix+"auth.credentials", "", "NATS credentials file")
	flags.StringVar(&c.Credentials.JWT, prefix+"auth.jwt", "", "NATS JWT file")
	flags.StringVar(&c.Credentials.Seed, prefix+"auth.seed", "", "NATS seed file")
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
