// Command protoget just gets your protos.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/flagenv"
)

var currentVersion = "0"

// Config is the configuration of the app.
type Config struct {
	Debug  bool
	Dst    string
	Prefix string
}

// App is the main app.
type App struct {
	config Config
	client *http.Client
}

var app = new(App)

func init() {
	flag.BoolVar(&app.config.Debug, "debug", false, "run in debug mode")
	flag.StringVar(&app.config.Dst, "dst", "", "destination dir")
	flag.StringVar(&app.config.Prefix, "prefix", "", "prefix to add before the path")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "usage: protoget [flags] [path ...]")
		flag.PrintDefaults()
	}
}

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	if err := flagenv.NewParser().ParseEnv(flag.CommandLine); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		flag.Usage()
		os.Exit(2)
	}

	flag.Parse()

	if app.config.Prefix != "" && !strings.HasSuffix(app.config.Prefix, "/") {
		app.config.Prefix += "/"
	}

	if err := app.Run(ctx, flag.Args()...); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		return
	}
}

// Run runs the app.
func (app *App) Run(ctx context.Context, args ...string) error {
	for _, arg := range args {
		if err := app.Get(ctx, arg); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (app *App) url(arg string) string {
	arg = app.config.Prefix + arg
	switch {
	case strings.HasPrefix(arg, "github.com/"), strings.HasPrefix(arg, "gitlab.com/"), strings.HasPrefix(arg, "bitbucket.org/"):
		parts := strings.SplitN(arg, "/", 4)
		if len(parts) == 4 {
			return fmt.Sprintf("https://%s/%s/%s/raw/master/%s", parts[0], parts[1], parts[2], parts[3])
		}
	case strings.HasPrefix(arg, "google/protobuf/"):
		return fmt.Sprintf("https://github.com/protocolbuffers/protobuf/raw/master/src/%s", arg)
	case strings.HasPrefix(arg, "google/"):
		return fmt.Sprintf("https://github.com/googleapis/googleapis/raw/master/%s", arg)
	}
	return "https://" + arg
}

func (app *App) Get(ctx context.Context, arg string) (err error) {
	url := app.url(arg)
	log.Printf("Getting %q from %q", arg, url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("protoget/%s", currentVersion))
	httpClient := app.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to get %q: %s", arg, res.Status)
	}
	err = os.MkdirAll(filepath.Dir(filepath.Join(app.config.Dst, arg)), 0755)
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(filepath.Join(app.config.Dst, arg), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dst.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	_, err = io.Copy(dst, res.Body)
	if err != nil {
		return err
	}
	return nil
}
