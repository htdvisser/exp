package main

import (
	"encoding/base32"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

var flags = func() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("b32", pflag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: b32 [FLAGS] [FILE]")
		fmt.Fprintln(os.Stderr, "Base32 encode or decode FILE, or standard input, to standard output")
		flagSet.PrintDefaults()
	}
	return flagSet
}()

var (
	decode = flags.BoolP("decode", "d", false, "decode data")
	raw    = flags.BoolP("raw", "r", false, "raw base32")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	encoding := base32.StdEncoding
	if *raw {
		encoding = encoding.WithPadding(base32.NoPadding)
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

	var w io.WriteCloser
	if *decode {
		r = base32.NewDecoder(encoding, r)
		w = os.Stdout
	} else {
		w = base32.NewEncoder(encoding, os.Stdout)
	}

	n, err := io.Copy(w, r)
	if err != nil {
		if n > 0 {
			fmt.Fprintln(os.Stdout)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = w.Close()
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
