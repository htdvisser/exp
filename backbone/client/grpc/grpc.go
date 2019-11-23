package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type dialOptionsContextKeyType struct{}

var dialOptionsContextKey dialOptionsContextKeyType

// NewContextWithDialOptions returns a context derived from parent that contains dialOptions.
func NewContextWithDialOptions(parent context.Context, dialOptions ...grpc.DialOption) context.Context {
	return context.WithValue(parent, dialOptionsContextKey, append(DialOptionsFromContext(parent), dialOptions...))
}

// DialOptionsFromContext returns the DialOptions from the context if present. Otherwise it returns nil.
func DialOptionsFromContext(ctx context.Context) []grpc.DialOption {
	if dialOptions, ok := ctx.Value(dialOptionsContextKey).([]grpc.DialOption); ok {
		options := make([]grpc.DialOption, len(dialOptions))
		copy(options, dialOptions)
		return options
	}
	return nil
}

// DialContext calls grpc.DialContext and includes the DialOptions in the context.
func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target, append(DialOptionsFromContext(ctx), opts...)...)
}

type callOptionsContextKeyType struct{}

var callOptionsContextKey callOptionsContextKeyType

// NewContextWithCallOptions returns a context derived from parent that contains callOptions.
func NewContextWithCallOptions(parent context.Context, callOptions ...grpc.CallOption) context.Context {
	return context.WithValue(parent, callOptionsContextKey, append(CallOptionsFromContext(parent), callOptions...))
}

// CallOptionsFromContext returns the CallOptions from the context if present. Otherwise it returns nil.
func CallOptionsFromContext(ctx context.Context) []grpc.CallOption {
	if callOptions, ok := ctx.Value(callOptionsContextKey).([]grpc.CallOption); ok {
		return callOptions
	}
	return nil
}
