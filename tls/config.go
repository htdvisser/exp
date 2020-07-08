package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
)

// FileReader reads a file and returns its bytes.
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// CertConfig represent a TLS certificate and its key.
type CertConfig struct {
	Certificate string
	Key         string
}

// Load loads the certificate. If fileReader is nil, this uses ioutil.ReadFile.
func (f CertConfig) Load(_ context.Context, fileReader FileReader) (*tls.Certificate, error) {
	if f.Certificate == "" && f.Key == "" {
		return nil, nil
	}
	readFile := ioutil.ReadFile
	if fileReader != nil {
		readFile = fileReader.ReadFile
	}
	certPEMBlock, err := readFile(f.Certificate)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := readFile(f.Key)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// CAConfig represents a CA certificate.
type CAConfig struct {
	Certificate string
}

// Load loads the CA certificates. If fileReader is nil, this uses ioutil.ReadFile.
func (f CAConfig) Load(_ context.Context, fileReader FileReader) ([]*x509.Certificate, error) {
	if f.Certificate == "" {
		return nil, nil
	}
	readFile := ioutil.ReadFile
	if fileReader != nil {
		readFile = fileReader.ReadFile
	}
	certPEMBlock, err := readFile(f.Certificate)
	if err != nil {
		return nil, err
	}
	var certs []*x509.Certificate
	for len(certPEMBlock) > 0 {
		var block *pem.Block
		block, certPEMBlock = pem.Decode(certPEMBlock)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// Config is the configuration for client and server TLS.
type Config struct {
	Store      string
	ClientCert CertConfig
	ClientCA   CAConfig
	ServerCert CertConfig
	ServerCA   CAConfig
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string) *pflag.FlagSet {
	var flags pflag.FlagSet
	flags.StringVar(&c.Store, prefix+"store", "file", "TLS Certificate store")
	flags.StringVar(&c.ClientCert.Certificate, prefix+"client.cert", "", "TLS Client certificate file")
	flags.StringVar(&c.ClientCert.Key, prefix+"client.key", "", "TLS Client certificate key file")
	flags.StringVar(&c.ClientCA.Certificate, prefix+"client.ca.cert", "", "TLS Client CA certificate file")
	flags.StringVar(&c.ServerCert.Certificate, prefix+"server.cert", "", "TLS Server certificate file")
	flags.StringVar(&c.ServerCert.Key, prefix+"server.key", "", "TLS Server certificate key file")
	flags.StringVar(&c.ServerCA.Certificate, prefix+"server.ca.cert", "", "TLS Server CA certificate file")
	return &flags
}

func (c Config) getStore() (FileReader, error) {
	switch c.Store {
	case "", "file":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported TLS certificate store %q", c.Store)
	}
}

// BuildServerConfig builds a TLS config suitable for use by servers.
func (c Config) BuildServerConfig(ctx context.Context) (*tls.Config, error) {
	var tlsConfig tls.Config
	tlsConfig.MinVersion = tls.VersionTLS12
	store, err := c.getStore()
	if err != nil {
		return nil, err
	}
	serverCert, err := c.ServerCert.Load(ctx, store)
	if err != nil {
		return nil, err
	}
	if serverCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*serverCert}
	}
	clientCAs, err := c.ClientCA.Load(ctx, store)
	if err != nil {
		return nil, err
	}
	if len(clientCAs) > 0 {
		pool := x509.NewCertPool()
		for _, clientCA := range clientCAs {
			pool.AddCert(clientCA)
		}
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
	}
	return &tlsConfig, nil
}

// BuildClientConfig builds a TLS config suitable for use by clients.
func (c Config) BuildClientConfig(ctx context.Context) (*tls.Config, error) {
	var tlsConfig tls.Config
	tlsConfig.MinVersion = tls.VersionTLS12
	store, err := c.getStore()
	if err != nil {
		return nil, err
	}
	clientCert, err := c.ClientCert.Load(ctx, store)
	if err != nil {
		return nil, err
	}
	if clientCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*clientCert}
	}
	serverCAs, err := c.ServerCA.Load(ctx, store)
	if err != nil {
		return nil, err
	}
	if len(serverCAs) > 0 {
		pool, err := x509.SystemCertPool()
		if err != nil {
			pool = x509.NewCertPool()
		}
		for _, serverCA := range serverCAs {
			pool.AddCert(serverCA)
		}
		tlsConfig.RootCAs = pool
	}
	return &tlsConfig, nil
}
