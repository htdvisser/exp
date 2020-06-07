package flagenv

import (
	"flag"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	emptyFlagSet := func() *flag.FlagSet {
		return flag.NewFlagSet("test", flag.ContinueOnError)
	}

	t.Run("Empty", func(t *testing.T) {
		p := NewParser(Prefixes("TEST_", ""), ReplaceWithUnderscore("-", "+"))
		err := p.ParseEnv(emptyFlagSet())
		if err != nil {
			t.Errorf("p.ParseEnv with empty FlagSet = %v, want nil", err)
		}
	})

	flagSet := func() *flag.FlagSet {
		flagSet := emptyFlagSet()
		flagSet.String("flag", "default flag value", "usage")
		return flagSet
	}

	for _, tt := range []struct {
		Name    string
		Options []ParserOption
		env     map[string]string
		flagSet func() *flag.FlagSet
		want    string
	}{
		{
			Name:    "Empty Env",
			env:     map[string]string{},
			flagSet: flagSet,
			want:    "default flag value",
		},
		{
			Name: "Prefixed Env",
			env: map[string]string{
				"TEST_FLAG": "test flag value",
			},
			flagSet: flagSet,
			want:    "test flag value",
		},
		{
			Name: "Non-Prefixed Env",
			env: map[string]string{
				"FLAG": "flag value",
			},
			flagSet: flagSet,
			want:    "flag value",
		},
		{
			Name: "Filtered Env",
			Options: []ParserOption{Filter(func(key string) bool {
				if key == "flag" {
					return false
				}
				return true
			})},
			env: map[string]string{
				"FLAG": "flag value",
			},
			flagSet: flagSet,
			want:    "default flag value",
		},
		{
			Name: "Priority Env",
			env: map[string]string{
				"TEST_FLAG": "test flag value",
				"FLAG":      "flag value",
			},
			flagSet: flagSet,
			want:    "flag value",
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			p := NewParser(append([]ParserOption{
				Prefixes("TEST_", ""), ReplaceWithUnderscore("-", "+"),
			}, tt.Options...)...)
			p.lookupEnv = func(key string) (string, bool) {
				val, ok := tt.env[key]
				return val, ok
			}

			flagSet := tt.flagSet()

			err := p.ParseEnv(flagSet)
			if err != nil {
				t.Errorf("p.ParseEnv with env (%v) = %v, want nil", tt.env, err)
			}

			if got := flagSet.Lookup("flag").Value.String(); got != tt.want {
				t.Errorf("output of flag %q = %q, want %q", "flag", got, tt.want)
			}
		})
	}
}

func TestParserError(t *testing.T) {
	p := NewParser(Prefixes("TEST_", ""), ReplaceWithUnderscore("-", "+"))
	env := map[string]string{
		"TEST_FLAG":  "error",
		"TEST_OTHER": "error",
	}
	p.lookupEnv = func(key string) (string, bool) {
		val, ok := env[key]
		return val, ok
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Int("flag", 0, "usage")
	flagSet.Int("other", 0, "usage")

	err := p.ParseEnv(flagSet)
	if err == nil {
		t.Errorf("p.ParseEnv with (%v) = nil, want error", env)
	}
	if !strings.Contains(err.Error(), "flagenv: invalid environment") {
		t.Errorf("error %q did not contain %q", err, "flagenv: invalid environment")
	}
}
