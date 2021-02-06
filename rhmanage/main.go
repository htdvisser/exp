package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/pflagenv"
	"htdvisser.dev/exp/redisconfig"
)

const bin = "rhmanage"

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
	flags       = pflag.NewFlagSet(bin, pflag.ContinueOnError)
	versionFlag = flags.BoolP("version", "V", false, "Print version information")
	redisConfig redisconfig.Config
	config      struct {
		Match  string
		Batch  int
		Filter map[string]string
		HSET   map[string]string
		HDEL   []string
		DEL    bool
		DryRun bool
	}
	logger = log.New(os.Stderr, "", log.LstdFlags)
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
	flags.AddFlagSet(redisConfig.Flags("redis.", nil))
	flags.StringVar(&config.Match, "key.match", "*", "MATCH the key")
	flags.IntVar(&config.Batch, "batch", 1000, "Batch size")
	flags.StringToStringVar(&config.Filter, "filter", nil, "Filter by hash contents")
	flags.StringToStringVar(&config.HSET, "hset", nil, "HMSET on filtered hashes")
	flags.StringSliceVar(&config.HDEL, "hdel", nil, "HDEL on filtered hashes")
	flags.BoolVar(&config.DEL, "del", false, "DEL the keys")
	flags.BoolVar(&config.DryRun, "dry-run", false, "Don't actually perform the action")
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
	redisCli, err := redisConfig.Connect(ctx)
	if err != nil {
		return err
	}
	defer redisCli.Close()

	filterFields := make([]string, 0, len(config.Filter))
	for field := range config.Filter {
		filterFields = append(filterFields, field)
	}
	hmset := make(map[string]interface{}, len(config.HSET))
	for k, v := range config.HSET {
		hmset[k] = v
	}

	var keys []string
	var cursor uint64

	for {
		logger.Printf("SCAN %d keys MATCH %s", config.Batch, config.Match)
		keys, cursor, err = redisCli.Scan(ctx, cursor, config.Match, int64(config.Batch)).Result()
		if err != nil {
			logger.Printf("SCAN error: %s", err)
			break
		}
		if len(keys) == 0 {
			logger.Println(ctx, "SCAN returned empty batch")
			continue
		}
		logger.Printf("SCAN returned %d keys", len(keys))

	NextKey:
		for _, key := range keys {
			if config.DEL {
				if !config.DryRun {
					err := redisCli.Del(ctx, key).Err()
					if err != nil {
						logger.Printf("DEL error: %s", err)
						continue
					}
				}
				logger.Printf("DEL %s", key)
				continue
			}

			keyType, err := redisCli.Type(ctx, key).Result()
			if err != nil {
				logger.Printf("TYPE error: %s", err)
				continue
			}
			if keyType != "hash" {
				logger.Printf("%s is not a hash", key)
				continue
			}

			if len(filterFields) > 0 {
				fieldValues, err := redisCli.HMGet(ctx, key, filterFields...).Result()
				if err != nil {
					logger.Printf("HMGET error: %s", err)
					continue
				}
				for i, value := range fieldValues {
					var matches bool
					switch value := value.(type) {
					case nil:
						matches = config.Filter[filterFields[i]] == ""
					case string:
						matches = value == config.Filter[filterFields[i]]
					default:
						logger.Printf("Don't know how to match %T", value)
					}
					if !matches {
						logger.Printf("%s does not match %s", key, filterFields[i])
						continue NextKey
					}
				}
			}

			if len(hmset) > 0 {
				if !config.DryRun {
					err := redisCli.HMSet(ctx, key, hmset).Err()
					if err != nil {
						logger.Printf("HMSET error: %s", err)
						continue
					}
				}
				logger.Printf("HMSET %s: %v", key, hmset)
			}

			if len(config.HDEL) > 0 {
				if !config.DryRun {
					err := redisCli.HDel(ctx, key, config.HDEL...).Err()
					if err != nil {
						logger.Printf("HDEL error: %s", err)
						continue
					}
				}
				logger.Printf("HDEL %s: %v", key, config.HDEL)
			}
		}

		if cursor == 0 {
			logger.Println("SCAN done")
			break
		}
	}

	return nil
}
