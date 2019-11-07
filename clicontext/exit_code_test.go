package clicontext

import (
	"context"
	"testing"
)

func TestExitCode(t *testing.T) {
	bgCtx := context.Background()

	t.Run("Background Context", func(t *testing.T) {
		code, ok := GetExitCode(bgCtx)
		if ok || code != 0 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 0, false)
		}

		ok = SetExitCode(bgCtx, 1)
		if ok {
			t.Errorf("SetExitCode(ctx, 1) = %v, want %v", ok, false)
		}

		code, ok = GetExitCode(bgCtx)
		if ok || code != 0 {
			t.Errorf("GetExitCode(ctx) after SetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 0, false)
		}
	})

	t.Run("Nil Destination", func(t *testing.T) {
		ctx := WithExitCode(bgCtx, nil)

		code, ok := GetExitCode(ctx)
		if !ok || code != 0 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 0, true)
		}

		ok = SetExitCode(ctx, 1)
		if !ok {
			t.Errorf("SetExitCode(ctx, 1) = %v, want %v", ok, true)
		}

		code, ok = GetExitCode(ctx)
		if !ok || code != 1 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 1, true)
		}
	})

	t.Run("With Destination", func(t *testing.T) {
		var code int

		ctx := WithExitCode(bgCtx, &code)

		ok := SetExitCode(ctx, 1)
		if !ok {
			t.Errorf("SetExitCode(ctx, 1) = %v, want %v", ok, true)
		}

		if code != 1 {
			t.Errorf("code after SetExitCode = %v, want %v", code, 1)
		}

		code, ok = GetExitCode(ctx)
		if !ok || code != 1 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 1, true)
		}
	})
}
