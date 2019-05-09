package flagenv

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := NewParser(Prefixes("TEST_", ""), ReplaceWithUnderscore("-", "+"))

	env := map[string]string{}
	p.lookupEnv = func(key string) (string, bool) {
		val, ok := env[key]
		return val, ok
	}

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	err := p.ParseEnv(flagSet)
	assert.NoError(t, err)

	out := flagSet.String("flag", "default flag value", "usage")

	err = p.ParseEnv(flagSet)
	assert.NoError(t, err)
	assert.Equal(t, "default flag value", *out)

	env = map[string]string{
		"TEST_FLAG": "test flag value",
	}

	err = p.ParseEnv(flagSet)
	assert.NoError(t, err)
	assert.Equal(t, "test flag value", *out)

	env = map[string]string{
		"FLAG": "flag value",
	}

	err = p.ParseEnv(flagSet)
	assert.NoError(t, err)
	assert.Equal(t, "flag value", *out)

	env = map[string]string{
		"TEST_FLAG": "test flag value",
		"FLAG":      "flag value",
	}

	err = p.ParseEnv(flagSet)
	assert.NoError(t, err)
	assert.Equal(t, "flag value", *out)
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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "flagenv: invalid environment")
}
