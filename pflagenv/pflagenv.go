// Package pflagenv helps with parsing environment into flags.
package pflagenv // import "htdvisser.dev/exp/pflagenv"

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// Parser parses environment into flags.
type Parser struct {
	lookupEnv func(key string) (string, bool)
	prefixes  []string
	replacer  *strings.Replacer
}

var (
	defaultReplacer = buildReplacer(".", "-")
	defaultPrefixes = []string{""}
)

func (p *Parser) setDefaults() {
	if p.lookupEnv == nil {
		p.lookupEnv = os.LookupEnv
	}
	if p.replacer == nil {
		p.replacer = defaultReplacer
	}
	if p.prefixes == nil {
		p.prefixes = defaultPrefixes
	}
}

// ParserOption configures a Parser.
type ParserOption interface {
	apply(*Parser)
}

type parserOptionFunc func(*Parser)

func (f parserOptionFunc) apply(p *Parser) { f(p) }

func buildReplacer(chars ...string) *strings.Replacer {
	oldnew := make([]string, 0, len(chars)*2)
	for _, char := range chars {
		oldnew = append(oldnew, char, "_")
	}
	return strings.NewReplacer(oldnew...)
}

// ReplaceWithUnderscore returns a ParserOption that makes the Parser replace
// the given characters with underscores.
func ReplaceWithUnderscore(chars ...string) ParserOption {
	return parserOptionFunc(func(p *Parser) {
		p.replacer = buildReplacer(chars...)
	})
}

// Prefixes returns a ParserOption that makes the Parser consider the given
// prefixes. It replaces configured characters in the prefix with underscores
// (see also ReplaceWithUnderscore). It does not add an underscore between the
// prefix and the flag name, so make sure to add an underscore if needed.
func Prefixes(prefixes ...string) ParserOption {
	return parserOptionFunc(func(p *Parser) {
		p.prefixes = make([]string, len(prefixes))
		for i, prefix := range prefixes {
			p.prefixes[i] = strings.ToUpper(p.replacer.Replace(prefix))
		}
	})
}

// NewParser returns a new Parser with the given options.
func NewParser(options ...ParserOption) *Parser {
	p := &Parser{}
	p.setDefaults()
	for _, option := range options {
		option.apply(p)
	}
	return p
}

func (p *Parser) parseEnv(flag *pflag.Flag) error {
	name := strings.ToUpper(p.replacer.Replace(flag.Name))
	for _, prefix := range p.prefixes {
		key := prefix + name
		val, present := p.lookupEnv(key)
		if !present {
			continue
		}
		if err := flag.Value.Set(val); err != nil {
			return fmt.Errorf("pflagenv: invalid environment %s=%q for flag -%s: %v", key, val, flag.Name, err)
		}
	}
	return nil
}

type errSlice []error

func (errs errSlice) Error() string {
	str := errs[0].Error()
	switch remaining := len(errs) - 1; remaining {
	case 0:
		return str
	case 1:
		return str + fmt.Sprintf(" (and %d more error)", remaining)
	default:
		return str + fmt.Sprintf(" (and %d more errors)", remaining)
	}
}

// ParseEnv parses the environment for the given FlagSet.
func (p *Parser) ParseEnv(flagSet *pflag.FlagSet) error {
	p.setDefaults()
	var errors []error
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if err := p.parseEnv(flag); err != nil {
			errors = append(errors, err)
		}
	})
	if len(errors) > 0 {
		return errSlice(errors)
	}
	return nil
}
