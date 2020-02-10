package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

var flags = func() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("hex", pflag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: hex [FLAGS] [FILE]")
		fmt.Fprintln(os.Stderr, "Hex encode or decode FILE, or standard input, to standard output")
		flagSet.PrintDefaults()
	}
	return flagSet
}()

var decode = flags.BoolP("decode", "d", false, "decode data")

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	args := flags.Args()
	var r io.Reader
	switch len(args) {
	case 0:
		r = os.Stdin
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	default:
		fmt.Fprintln(os.Stderr, "invalid number of arguments")
		os.Exit(2)
	}

	var w io.Writer
	if *decode {
		r = hex.NewDecoder(&newlineFilteringReader{r})
		w = os.Stdout
	} else {
		w = hex.NewEncoder(os.Stdout)
	}

	n, err := io.Copy(w, r)
	if err != nil {
		if n > 0 {
			fmt.Fprintln(os.Stdout)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !*decode && n > 0 {
		fmt.Fprintln(os.Stdout)
	}
}
