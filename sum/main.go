package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var flags = func() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("sum", pflag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sum [TYPE] [FILE]")
		fmt.Fprintln(os.Stderr, "Sum checksums FILE, or standard input, to standard output")
		flagSet.PrintDefaults()
	}
	return flagSet
}()

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
	var (
		r    io.Reader
		hash hash.Hash
	)
	switch len(args) {
	case 0:
		flags.Usage()
		os.Exit(2)
	case 1:
		r = os.Stdin
		hash, err = selectHash(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	case 2:
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
		hash, err = selectHash(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid number of arguments")
		os.Exit(2)
	}

	_, err = io.Copy(hash, r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	n, err := os.Stdout.Write(hash.Sum(nil))
	if err != nil {
		if n > 0 {
			fmt.Fprintln(os.Stdout)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func selectHash(kind string) (hash.Hash, error) {
	switch strings.ToLower(kind) {
	case "crc32":
		return crc32.New(crc32.MakeTable(crc32.IEEE)), nil
	case "crc64":
		return crc64.New(crc64.MakeTable(crc64.ISO)), nil
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash: %q", kind)
	}
}
