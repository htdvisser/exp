package clicontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCode(t *testing.T) {
	bgCtx := context.Background()

	t.Run("Background Context", func(t *testing.T) {
		code, ok := GetExitCode(bgCtx)
		assert.False(t, ok)
		assert.Zero(t, code)

		ok = SetExitCode(bgCtx, 1)
		assert.False(t, ok)

		code, ok = GetExitCode(bgCtx)
		assert.False(t, ok)
		assert.Zero(t, code)
	})

	t.Run("Nil Destination", func(t *testing.T) {
		ctx := WithExitCode(bgCtx, nil)

		code, ok := GetExitCode(ctx)
		assert.True(t, ok)
		assert.Equal(t, 0, code)

		ok = SetExitCode(ctx, 1)
		assert.True(t, ok)

		code, ok = GetExitCode(ctx)
		assert.True(t, ok)
		assert.Equal(t, 1, code)
	})

	t.Run("With Destination", func(t *testing.T) {
		var code int

		ctx := WithExitCode(bgCtx, &code)

		ok := SetExitCode(ctx, 1)
		assert.True(t, ok)

		assert.Equal(t, 1, code)
	})
}
