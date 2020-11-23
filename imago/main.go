package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/pflag"
)

var flags = func() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("imago", pflag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: imago [FLAGS] [FILE...]")
		fmt.Fprintln(os.Stderr, "A tool to resize and convert images")
		flagSet.PrintDefaults()
	}
	return flagSet
}()

var (
	out     = flags.StringP("out", "o", ".", "output dir")
	quality = flags.IntP("quality", "q", 70, "JPEG quality")
	jpeg    = flags.Bool("jpeg", true, "export as JPEG")
	fit     = flags.Int("fit", 0, "fit to size (default disabled)")
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

	for _, arg := range flags.Args() {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %q: %v\n", arg, err)
			os.Exit(1)
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode %q: %v\n", arg, err)
			os.Exit(1)
		}
		if bounds := img.Bounds(); *fit > 0 && (bounds.Dy() > *fit || bounds.Dx() > *fit) {
			img = imaging.Fit(img, *fit, *fit, imaging.Lanczos)
		}
		outfilename := filepath.Base(arg)
		if ext := filepath.Ext(outfilename); *jpeg && ext != ".jpg" {
			outfilename = strings.TrimSuffix(outfilename, ext) + ".jpg"
		}
		err = imaging.Save(img, filepath.Join(*out, outfilename), imaging.JPEGQuality(*quality))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode %q to %q: %v\n", arg, outfilename, err)
			os.Exit(1)
		}
	}
}
