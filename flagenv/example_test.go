package flagenv_test

import (
	"flag"
	"fmt"
	"os"

	"htdvisser.dev/exp/flagenv"
)

func Example() {
	flag.Bool("debug", false, "enable debug mode")

	// Environment: DEBUG=1

	if err := flagenv.NewParser().ParseEnv(flag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	flag.Parse()
}

func ExamplePrefixes() {
	flag.Bool("debug", false, "enable debug mode")

	// Environment: FOO_DEBUG=1

	if err := flagenv.NewParser(flagenv.Prefixes("FOO_")).ParseEnv(flag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	flag.Parse()
}
