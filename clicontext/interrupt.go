package clicontext

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var defaultInterruptSignals = []os.Signal{os.Interrupt}

// WithInterrupt returns a copy of parent with a new Done channel, which is closed
// when the process receives an interrupt signal or when the parent context's Done
// channel is closed.
// The optional extra interrupt signals are added to the default interrupt signals.
func WithInterrupt(parent context.Context, extraInterruptSignals ...os.Signal) context.Context {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, append(defaultInterruptSignals[:], extraInterruptSignals...)...)
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-parent.Done():
		case osSig := <-sig:
			if sig, ok := osSig.(syscall.Signal); ok {
				SetExitCode(ctx, int(sig)+128)
			}
			cancel()
		}
		signal.Stop(sig)
	}()
	return ctx
}

// WithTemporaryInterrupt is similar to WithInterrupt, but also releases resources
// when the returned context.CancelFunc is called.
func WithTemporaryInterrupt(parent context.Context, extraInterruptSignals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return WithInterrupt(ctx, extraInterruptSignals...), cancel
}

// WithInterruptAndExit is similar to WithInterrupt, but also returns a func that
// exits the program with an appropriate status code for the interrupt signal.
func WithInterruptAndExit(parent context.Context, extraInterruptSignals ...os.Signal) (ctx context.Context, exit func()) {
	var code int
	ctx = WithInterrupt(WithExitCode(parent, &code), extraInterruptSignals...)
	return ctx, func() { os.Exit(code) }
}
