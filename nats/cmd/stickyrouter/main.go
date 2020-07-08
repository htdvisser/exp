// Command stickyrouter subscribes to a NATS subject in the form of
// `sticky.route.{duration}.{hash}`. From the first time a message is received,
// on a subject, all messages for the same `{hash}` are routed to the same
// reply subject for a duration of `{duration}`.
//
// The `{duration}` is given as a string in the form of 1h2m3.45s.
package main

import (
	"context"
	"expvar" // Registers /debug/vars to http.DefaultServeMux.
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof" // Registers /debug/pprof to http.DefaultServeMux.
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/nats"
	"htdvisser.dev/exp/nats/internal/stickyrouter"
	"htdvisser.dev/exp/pflagenv"
	"htdvisser.dev/exp/redis"
	"htdvisser.dev/exp/tls"
)

const bin = "stickyrouter"

// Variables set during build (using ldflags):
var (
	version = "0.0.0" // -X main.version=0.0.0
	commit  = ""      // -X main.commit=$(git rev-parse HEAD)
	date    = ""      // -X main.date=$(date -uIseconds)
)

func init() {
	version = strings.TrimPrefix(version, "v")
	expvar.NewString("version").Set(version)
	expvar.NewString("commit").Set(commit)
	expvar.NewString("builddate").Set(date)
}

var (
	flags              = pflag.NewFlagSet(bin, pflag.ContinueOnError)
	versionFlag        = flags.BoolP("version", "V", false, "Print version information")
	debugAddrFlag      = flags.String("debug.addr", "localhost:6060", "Address to listen on for debug endpoints")
	tlsConfig          tls.Config
	natsConfig         nats.Config
	redisConfig        redis.Config
	stickyrouterConfig stickyrouter.Config
)

func filterEnv(key string) bool {
	switch key {
	case "version":
		return false
	default:
		return true
	}
}

func init() {
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", bin)
		flags.PrintDefaults()
	}
	flags.AddFlagSet(tlsConfig.Flags("tls."))
	flags.AddFlagSet(natsConfig.Flags("nats."))
	flags.AddFlagSet(redisConfig.Flags("redis."))
	flags.AddFlagSet(stickyrouterConfig.Flags("route."))
}

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	defer func() {
		if p := recover(); p != nil {
			clicontext.SetExitCode(ctx, 1)
			panic(p)
		}
	}()

	if err := pflagenv.NewParser(
		pflagenv.Filter(filterEnv),
	).ParseEnv(flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		clicontext.SetExitCode(ctx, 2)
		return
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		if err != pflag.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			flags.Usage()
		}
		clicontext.SetExitCode(ctx, 2)
		return
	}

	if *versionFlag || len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Fprintf(
			os.Stdout,
			"%s %s %s %s/%s\n",
			bin, version, runtime.Version(), runtime.GOOS, runtime.GOARCH,
		)
		clicontext.SetExitCode(ctx, 0)
		return
	}

	if err := Run(ctx, os.Args[1:]...); err != nil {
		if err != context.Canceled {
			clicontext.SetExitCode(ctx, 1)
		}
		fmt.Fprintln(os.Stderr, err)
	}
}

// Run runs the program until the Done() channel of the context is closed.
func Run(ctx context.Context, args ...string) error {
	natsConn, err := natsConfig.Connect(ctx)
	if err != nil {
		return err
	}
	defer natsConn.Close()

	redisCli, err := redisConfig.Connect(ctx)
	if err != nil {
		return err
	}
	defer redisCli.Close()

	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	srv := http.Server{Addr: *debugAddrFlag}
	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	log.Printf("Listen on %q for debug server", lis.Addr())
	errGroup.Go(func() error {
		if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	errGroup.Go(func() error {
		<-errGroupCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	svc, err := stickyrouter.NewService(&stickyrouterConfig, natsConn, redisCli)
	if err != nil {
		return err
	}
	errGroup.Go(func() error {
		return svc.Run(errGroupCtx)
	})

	return errGroup.Wait()
}
