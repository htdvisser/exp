package tlsconfig

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/pflag"
)

func TestCertConfigFlags(t *testing.T) {
	var config CertConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.cert=cert.pem",
		"test.key=cert-key.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}

func TestCAConfigFlags(t *testing.T) {
	var config CAConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.caCert=ca.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}

func TestServerConfigFlags(t *testing.T) {
	var config ServerConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.server.cert=server.pem",
		"test.server.key=server-key.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}

func TestMutualServerConfigFlags(t *testing.T) {
	var config MutualServerConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.server.cert=server.pem",
		"test.server.key=server-key.pem",
		"test.client.caCert=client-ca.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}

func TestClientConfigFlags(t *testing.T) {
	var config ClientConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.server.caCert=ca.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}

func TestMutualClientConfigFlags(t *testing.T) {
	var config MutualClientConfig
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.server.caCert=ca.pem",
		"test.client.cert=client.pem",
		"test.client.key=client-key.pem",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}
