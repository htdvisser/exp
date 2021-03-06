package redis

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/pflag"
)

// Config is the configuration for connecting to Redis.
type Config struct {
	Addresses    []string
	Username     string
	Password     string
	PasswordFile string
	PoolSize     int
	TLSConfig    *tls.Config
}

// DefaultConfig returns the default Redis configuration.
func DefaultConfig() *Config {
	return &Config{
		Addresses: []string{"localhost:6379"},
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string, defaults *Config) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultConfig()
	}
	flags.StringSliceVar(&c.Addresses, prefix+"addresses", defaults.Addresses, "Redis addresses")
	flags.StringVar(&c.Username, prefix+"auth.username", defaults.Username, "Redis username")
	flags.StringVar(&c.Password, prefix+"auth.password", defaults.Password, "Redis password")
	flags.StringVar(&c.PasswordFile, prefix+"auth.password-file", defaults.PasswordFile, "Redis password file")
	flags.IntVar(&c.PoolSize, prefix+"pool.size", defaults.PoolSize, "Redis connection pool size")
	return &flags
}

// Connect connects to Redis using this configuration.
func (c Config) Connect(ctx context.Context) (redis.UniversalClient, error) {
	opts := redis.UniversalOptions{
		Addrs:     c.Addresses,
		Username:  c.Username,
		Password:  c.Password,
		PoolSize:  c.PoolSize,
		TLSConfig: c.TLSConfig,
	}
	if c.PasswordFile != "" {
		passwordBytes, err := ioutil.ReadFile(c.PasswordFile)
		if err != nil {
			return nil, err
		}
		opts.Password = string(bytes.TrimSpace(passwordBytes))
	}
	cli := redis.NewUniversalClient(&opts)
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return cli, nil
}
