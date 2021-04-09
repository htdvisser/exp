package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
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
		Match        string
		Batch        int
		Filter       map[string]string
		FilterPrefix map[string]string
		HSET         map[string]string
		HDEL         []string
		DEL          bool
		GroupCount   string
		DryRun       bool
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
	flags.StringToStringVar(&config.FilterPrefix, "filter-prefix", nil, "Filter by prefix in hash contents")
	flags.StringToStringVar(&config.HSET, "hset", nil, "HMSET on filtered hashes")
	flags.StringSliceVar(&config.HDEL, "hdel", nil, "HDEL on filtered hashes")
	flags.StringVar(&config.GroupCount, "group-count", "", "Group by hash contents and count")
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

	var fieldsOfInterest []string
	filterFields := make([]string, 0, len(config.Filter))
	for field := range config.Filter {
		filterFields = append(filterFields, field)
	}
	fieldsOfInterest = append(fieldsOfInterest, filterFields...)
	filterPrefixFields := make([]string, 0, len(config.FilterPrefix))
	for field := range config.FilterPrefix {
		filterPrefixFields = append(filterPrefixFields, field)
	}
	fieldsOfInterest = append(fieldsOfInterest, filterPrefixFields...)
	if config.GroupCount != "" {
		fieldsOfInterest = append(fieldsOfInterest, config.GroupCount)
	}

	hmset := make(map[string]interface{}, len(config.HSET))
	for k, v := range config.HSET {
		hmset[k] = v
	}

	var (
		keys       []string
		cursor     uint64
		count      int
		groupCount = make(map[string]int)
	)

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

			if len(fieldsOfInterest) > 0 {
				fieldValues, err := redisCli.HMGet(ctx, key, fieldsOfInterest...).Result()
				if err != nil {
					logger.Printf("HMGET error: %s", err)
					continue
				}
				hash := make(map[string]interface{})
				for i, f := range fieldsOfInterest {
					hash[f] = fieldValues[i]
				}

				for k, filter := range config.Filter {
					var matches bool
					switch value := hash[k].(type) {
					case nil:
						matches = filter == ""
					case string:
						matches = value == filter
					default:
						logger.Printf("Don't know how to match %T", value)
					}
					if !matches {
						continue NextKey
					}
				}
				for k, filter := range config.FilterPrefix {
					var matches bool
					switch value := hash[k].(type) {
					case nil:
						matches = filter == ""
					case string:
						matches = strings.HasPrefix(value, filter)
					default:
						logger.Printf("Don't know how to match %T", value)
					}
					if !matches {
						continue NextKey
					}
				}
				if config.GroupCount != "" {
					switch value := hash[config.GroupCount].(type) {
					case string:
						groupCount[value]++
					default:
						logger.Printf("Can't group-count %T", value)
					}
				}
			}

			count++

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

	logger.Printf("Matched %d keys", count)

	if len(groupCount) > 0 {
		logger.Print("group count:")
		keys := make([]string, 0, len(groupCount))
		for k := range groupCount {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%d\t%s\n", groupCount[k], k)
		}
	}

	return nil
}
