package clicontext

import "context"

type exitCodeKeyType struct{}

var exitCodeKey exitCodeKeyType

// WithExitCode returns a context that can be used to set exit code
func WithExitCode(parent context.Context, dst *int) context.Context {
	if dst == nil {
		var code int
		dst = &code
	}
	return context.WithValue(parent, exitCodeKey, dst)
}

// GetExitCode returns the exit code from the context.
func GetExitCode(ctx context.Context) (code int, ok bool) {
	if dst, ok := ctx.Value(exitCodeKey).(*int); ok {
		return *dst, true
	}
	return 0, false
}

// SetExitCode sets the exit code into the context and returns true if it did so
// successfully.
func SetExitCode(ctx context.Context, code int) bool {
	if dst, ok := ctx.Value(exitCodeKey).(*int); ok {
		*dst = code
		return true
	}
	return false
}
