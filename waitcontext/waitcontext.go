// Package waitcontext implements a mechanism for making a context's CancelFunc
// wait for goroutines to finish.
package waitcontext

import (
	"context"
	"sync"
)

// CancelAndWaitFunc cancels the context and waits for all
type CancelAndWaitFunc func()

type waitContextKeyType struct{}

var waitContextKey waitContextKeyType

type waitContext struct {
	wg     sync.WaitGroup
	parent *waitContext
}

func (ctx *waitContext) add() {
	ctx.wg.Add(1)
	if ctx.parent != nil {
		ctx.parent.add()
	}
}

func (ctx *waitContext) done() {
	ctx.wg.Done()
	if ctx.parent != nil {
		ctx.parent.done()
	}
}

func (ctx *waitContext) wait() {
	ctx.wg.Wait()
}

// New returns a new context derived from parent, and a function that cancels the
// context and waits until all goroutines started with waitcontext.Go have finished.
func New(parent context.Context) (context.Context, CancelAndWaitFunc) {
	var wait waitContext
	if parentWait, ok := parent.Value(waitContextKey).(*waitContext); ok {
		wait.parent = parentWait
	}
	ctx, cancel := context.WithCancel(context.WithValue(parent, waitContextKey, &wait))
	return ctx, func() {
		cancel()
		wait.wait()
	}
}

// Go starts a new goroutine for f, and makes any CancelAndWaitFunc created for
// (a parent of) ctx wait for it to finish.
func Go(ctx context.Context, f func()) {
	done := noop
	if wait, ok := ctx.Value(waitContextKey).(*waitContext); ok {
		wait.add()
		done = wait.done
	}
	go func() {
		defer done()
		f()
	}()
}

func noop() {}
