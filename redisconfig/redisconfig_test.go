package redisconfig

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
		"test.addresses=[localhost:6379]",
		"test.username=",
		"test.password=",
		"test.passwordFile=",
		"test.poolSize=0",
	}); diff != "" {
		t.Errorf("Flags not as expected: %v", diff)
	}
}
