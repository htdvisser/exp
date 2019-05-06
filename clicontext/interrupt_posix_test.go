//+build !windows

package clicontext

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithInterrupt(t *testing.T) {
	ctx := WithInterrupt(context.Background(), syscall.SIGUSR1)

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	process.Signal(syscall.SIGUSR1)

	select {
	case <-ctx.Done():
		assert.Error(t, ctx.Err())
	case <-time.After(100 * time.Millisecond):
		t.Fatal("signal not received")
	}
}
