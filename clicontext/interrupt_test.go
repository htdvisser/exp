package clicontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTemporaryInterrupt(t *testing.T) {
	ctx, cancel := WithTemporaryInterrupt(context.Background())
	cancel()
	assert.Error(t, ctx.Err())
}
