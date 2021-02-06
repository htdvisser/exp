package stickyrouter

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// Config is the configuration for the sticky router service.
type Config struct {
	SubjectPattern string
	Queue          string
	Workers        int

	subject          string
	subjectTokens    int
	durationTokenIdx int
	hashTokenIdx     int
}

// DefaultConfig is the default configuration for the sticky router service.
func DefaultConfig() *Config {
	return &Config{
		SubjectPattern: "sticky.route.{duration}.{hash}",
		Queue:          "default",
		Workers:        1,
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string, defaults *Config) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultConfig()
	}
	flags.StringVar(&c.SubjectPattern, prefix+"subject", defaults.SubjectPattern, "Subject pattern to subscribe to")
	flags.StringVar(&c.Queue, prefix+"queue", defaults.Queue, "Queue to use when subscribing")
	flags.IntVar(&c.Workers, prefix+"workers", defaults.Workers, "Number of workers")
	return &flags
}

func (c *Config) parseSubjectPattern() error {
	if strings.ContainsAny(c.SubjectPattern, " \t\r\n") {
		return fmt.Errorf("subject pattern %q contains whitespace", c.SubjectPattern)
	}
	tokens := strings.Split(c.SubjectPattern, ".")
	c.subjectTokens = len(tokens)
	var (
		subjectTokens = make([]string, 0, len(tokens))
		gotDuration   bool
		gotHash       bool
	)
	for i, token := range tokens {
		switch token {
		case "{duration}":
			if gotDuration {
				return fmt.Errorf("subject pattern %q contains more than one {duration}", c.SubjectPattern)
			}
			gotDuration = true
			c.durationTokenIdx = i
			subjectTokens = append(subjectTokens, "*")
		case "{hash}":
			if gotHash {
				return fmt.Errorf("subject pattern %q contains more than one {hash}", c.SubjectPattern)
			}
			gotHash = true
			c.hashTokenIdx = i
			subjectTokens = append(subjectTokens, "*")
		default:
			subjectTokens = append(subjectTokens, token)
		}
	}
	if !gotDuration {
		return fmt.Errorf("subject pattern %q does not contain {duration}", c.SubjectPattern)
	}
	if !gotHash {
		return fmt.Errorf("subject pattern %q does not contain {hash}", c.SubjectPattern)
	}
	c.subject = strings.Join(subjectTokens, ".")
	return nil
}
