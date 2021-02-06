// Package tlsconfig provides config and flags for building TLS configurations for servers and clients.
package tlsconfig

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
)

// CertConfig represent a TLS certificate and its key.
type CertConfig struct {
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}

// DefaultCertConfig returns the default CertConfig.
func DefaultCertConfig() *CertConfig {
	return &CertConfig{
		Cert: "cert.pem",
		Key:  "cert-key.pem",
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *CertConfig) Flags(prefix string, defaults *CertConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultCertConfig()
	}
	flags.StringVar(&c.Cert, prefix+"cert", defaults.Cert, "TLS certificate file")
	flags.StringVar(&c.Key, prefix+"key", defaults.Key, "TLS certificate key file")
	return &flags
}

// Load loads the certificate.
func (c *CertConfig) Load(_ context.Context) (*tls.Certificate, error) {
	if c.Cert == "" && c.Key == "" {
		return nil, nil
	}
	certPEMBlock, err := ioutil.ReadFile(c.Cert)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate file: %w", err)
	}
	keyPEMBlock, err := ioutil.ReadFile(c.Key)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate key file: %w", err)
	}
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate key pair: %w", err)
	}
	return &cert, nil
}

// CAConfig represents a CA certificate.
type CAConfig struct {
	CACert string `json:"caCert,omitempty" yaml:"caCert,omitempty"`
}

// DefaultCAConfig returns the default CAConfig.
func DefaultCAConfig() *CAConfig {
	return &CAConfig{
		CACert: "ca.pem",
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *CAConfig) Flags(prefix string, defaults *CAConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultCAConfig()
	}
	flags.StringVar(&c.CACert, prefix+"caCert", defaults.CACert, "CA certificate file")
	return &flags
}

// Load loads the CA certificates.
func (c *CAConfig) Load(_ context.Context) ([]*x509.Certificate, error) {
	if c.CACert == "" {
		return nil, nil
	}
	certPEMBlock, err := ioutil.ReadFile(c.CACert)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate file: %w", err)
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

// ServerConfig is the configuration for server-side TLS.
type ServerConfig struct {
	ServerCert CertConfig `json:"server,omitempty" yaml:"server,omitempty"`
}

// DefaultServerConfig returns the default configuration.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		ServerCert: CertConfig{
			Cert: "server.pem",
			Key:  "server-key.pem",
		},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *ServerConfig) Flags(prefix string, defaults *ServerConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultServerConfig()
	}
	flags.AddFlagSet(c.ServerCert.Flags(prefix+"server.", &defaults.ServerCert))
	return &flags
}

// Load loads the TLS config.
func (c *ServerConfig) Load(ctx context.Context) (*tls.Config, error) {
	var tlsConfig tls.Config
	tlsConfig.MinVersion = tls.VersionTLS12
	serverCert, err := c.ServerCert.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not load server certificate: %w", err)
	}
	if serverCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*serverCert}
	}
	return &tlsConfig, nil
}

// MutualServerConfig is the configuration for server-side mTLS.
type MutualServerConfig struct {
	ServerConfig
	ClientCA CAConfig `json:"client,omitempty" yaml:"client,omitempty"`
}

// DefaultMutualServerConfig returns the default configuration.
func DefaultMutualServerConfig() *MutualServerConfig {
	return &MutualServerConfig{
		ServerConfig: *DefaultServerConfig(),
		ClientCA: CAConfig{
			CACert: "client-ca.pem",
		},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *MutualServerConfig) Flags(prefix string, defaults *MutualServerConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultMutualServerConfig()
	}
	flags.AddFlagSet(c.ServerCert.Flags(prefix+"server.", &defaults.ServerCert))
	flags.AddFlagSet(c.ClientCA.Flags(prefix+"client.", &defaults.ClientCA))
	return &flags
}

// Load loads the TLS config.
func (c *MutualServerConfig) Load(ctx context.Context) (*tls.Config, error) {
	tlsConfig, err := c.ServerConfig.Load(ctx)
	if err != nil {
		return nil, err // ServerConfig.Load wraps errors.
	}
	clientCAs, err := c.ClientCA.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not load client CA: %w", err)
	}
	if len(clientCAs) > 0 {
		pool := x509.NewCertPool()
		for _, clientCA := range clientCAs {
			pool.AddCert(clientCA)
		}
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
	}
	return tlsConfig, nil
}

// ClientConfig is the configuration for client-side TLS.
type ClientConfig struct {
	ServerCA CAConfig `json:"server,omitempty" yaml:"server,omitempty"`
}

// DefaultClientConfig returns the default configuration.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerCA: CAConfig{
			CACert: "ca.pem",
		},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *ClientConfig) Flags(prefix string, defaults *ClientConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultClientConfig()
	}
	flags.AddFlagSet(c.ServerCA.Flags(prefix+"server.", &defaults.ServerCA))
	return &flags
}

// Load loads the TLS config.
func (c *ClientConfig) Load(ctx context.Context) (*tls.Config, error) {
	var tlsConfig tls.Config
	tlsConfig.MinVersion = tls.VersionTLS12
	serverCAs, err := c.ServerCA.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not load server CA: %w", err)
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

// MutualClientConfig is the configuration for client-side mTLS.
type MutualClientConfig struct {
	ClientConfig
	ClientCert CertConfig `json:"client,omitempty" yaml:"client,omitempty"`
}

// DefaultMutualClientConfig returns the default configuration.
func DefaultMutualClientConfig() *MutualClientConfig {
	return &MutualClientConfig{
		ClientConfig: *DefaultClientConfig(),
		ClientCert: CertConfig{
			Cert: "client.pem",
			Key:  "client-key.pem",
		},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *MutualClientConfig) Flags(prefix string, defaults *MutualClientConfig) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultMutualClientConfig()
	}
	flags.AddFlagSet(c.ServerCA.Flags(prefix+"server.", &defaults.ServerCA))
	flags.AddFlagSet(c.ClientCert.Flags(prefix+"client.", &defaults.ClientCert))
	return &flags
}

// Load loads the TLS config.
func (c *MutualClientConfig) Load(ctx context.Context) (*tls.Config, error) {
	tlsConfig, err := c.ClientConfig.Load(ctx)
	if err != nil {
		return nil, err // ClientConfig.Load wraps errors.
	}
	clientCert, err := c.ClientCert.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not load client certificate: %w", err)
	}
	if clientCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*clientCert}
	}
	return tlsConfig, nil
}
