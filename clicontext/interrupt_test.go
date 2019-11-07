//+build !windows

package clicontext

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWithInterrupt(t *testing.T) {
	ctx := WithExitCode(context.Background(), nil)

	ctx = WithInterrupt(ctx, syscall.SIGUSR1)

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process for current pid: %v", err)
	}

	process.Signal(syscall.SIGUSR1)

	select {
	case <-ctx.Done():
		code, ok := GetExitCode(ctx)
		if !ok || code != 158 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 158, true)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("signal not received")
	}
}

func TestWithTemporaryInterrupt(t *testing.T) {
	ctx := WithExitCode(context.Background(), nil)

	ctx, cancel := WithTemporaryInterrupt(ctx, syscall.SIGUSR1)
	cancel()

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process for current pid: %v", err)
	}

	process.Signal(syscall.SIGUSR1)

	select {
	case <-ctx.Done():
		code, ok := GetExitCode(ctx)
		if !ok || code != 0 {
			t.Errorf("GetExitCode(ctx) = (%v, %v), want (%v, %v)", code, ok, 0, true)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("signal not received")
	}
}
