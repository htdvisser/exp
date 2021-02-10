package natsconfig

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/pflag"
)

func TestConfigFlags(t *testing.T) {
	var config Config
	flags := config.Flags("test.", nil)

	var flagNamesAndValues []string
	flags.VisitAll(func(f *pflag.Flag) {
		flagNamesAndValues = append(flagNamesAndValues, fmt.Sprintf("%s=%s", f.Name, f.DefValue))
	})

	if diff := cmp.Diff(flagNamesAndValues, []string{
		"test.servers=[nats://127.0.0.1:4222]",
		"test.name=",
		"test.auth.username=",
		"test.auth.password=",
		"test.auth.passwordFile=",
		"test.auth.credentialsFile=",
		"test.auth.jwtFile=",
		"test.auth.seedFile=",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}
