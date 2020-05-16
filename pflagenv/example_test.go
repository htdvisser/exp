package pflagenv_test

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"htdvisser.dev/exp/pflagenv"
)

func Example() {
	pflag.Bool("debug", false, "enable debug mode")

	// Environment: DEBUG=1

	if err := pflagenv.NewParser().ParseEnv(pflag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		pflag.Usage()
		os.Exit(2)
	}

	pflag.Parse()
}

func ExamplePrefixes() {
	pflag.Bool("debug", false, "enable debug mode")

	// Environment: FOO_DEBUG=1

	if err := pflagenv.NewParser(pflagenv.Prefixes("FOO_")).ParseEnv(pflag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		pflag.Usage()
		os.Exit(2)
	}

	pflag.Parse()
}
