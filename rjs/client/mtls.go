package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewMTLSClient returns a new Client that uses TLS Client Authentication to connect to the RJS server.
// It loads the CA, certificate and certificate key from the given files.
// Same as New it prepends the given base URL to API URIs, and sets the given Authorization header if not empty.
func NewMTLSClient(baseURL, caFile, certFile, keyFile, authorization string) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	if caFile != "" {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate from %q: %w", caFile, err)
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)
		transport.TLSClientConfig.RootCAs = certPool
	}
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read client certificate from %q and %q: %w", certFile, keyFile, err)
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
	if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport.Proxy = defaultTransport.Proxy
		transport.DialContext = defaultTransport.DialContext
		transport.MaxIdleConns = defaultTransport.MaxIdleConns
		transport.IdleConnTimeout = defaultTransport.IdleConnTimeout
		transport.TLSHandshakeTimeout = defaultTransport.TLSHandshakeTimeout
	}
	return New(&http.Client{
		Transport: transport,
	}, baseURL, authorization)
}
