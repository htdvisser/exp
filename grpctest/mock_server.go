package grpctest

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	testpb "google.golang.org/grpc/interop/grpc_testing"
)

// NewMockServer returns a new MockServer.
func NewMockServer(serverOptions []grpc.ServerOption, dialOptions []grpc.DialOption) *MockServer {
	return &MockServer{
		serverOptions: serverOptions,
		dialOptions:   dialOptions,
	}
}

// MockServer is a mock gRPC server that exposes the "grpc.testing" service.
type MockServer struct {
	serverOptions []grpc.ServerOption
	dialOptions   []grpc.DialOption

	testpb.UnimplementedTestServiceServer

	// Test call handlers.
	emptyCall           func(context.Context, *testpb.Empty) (*testpb.Empty, error)
	unaryCall           func(context.Context, *testpb.SimpleRequest) (*testpb.SimpleResponse, error)
	streamingOutputCall func(*testpb.StreamingOutputCallRequest, testpb.TestService_StreamingOutputCallServer) error
	streamingInputCall  func(testpb.TestService_StreamingInputCallServer) error
	fullDuplexCall      func(testpb.TestService_FullDuplexCallServer) error
	halfDuplexCall      func(testpb.TestService_HalfDuplexCallServer) error

	// Set by Start().
	lis    net.Listener
	server *grpc.Server

	// Called by Stop().
	cleanups []func()
}

// Start starts the mock server.
// The caller is responsible for stopping the server when done.
func (ms *MockServer) Start(ctx context.Context) error {
	if ms.lis != nil {
		return fmt.Errorf("mock server already started")
	}
	var err error
	ms.lis, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		return fmt.Errorf("failed to create listener for mock server: %w", err)
	}
	ms.cleanups = append(ms.cleanups, func() {
		ms.lis.Close()
		ms.lis = nil
	})
	ms.server = grpc.NewServer(ms.serverOptions...)
	testpb.RegisterTestServiceServer(ms.server, ms)
	go ms.server.Serve(ms.lis)
	ms.cleanups = append(ms.cleanups, func() {
		ms.server.Stop()
		ms.server = nil
	})
	return nil
}

// Stop stops the mock server.
func (ms *MockServer) Stop() {
	for i := len(ms.cleanups) - 1; i >= 0; i-- {
		ms.cleanups[i]()
	}
	ms.cleanups = nil
}

// ConnectTestClient connects to the mock server and returns a test client.
// The caller is responsible for closing the ClientConn when done.
func (ms *MockServer) ConnectTestClient(ctx context.Context) (*grpc.ClientConn, testpb.TestServiceClient, error) {
	if ms.lis == nil {
		return nil, nil, fmt.Errorf("mock server needs to be started before connecting test clients")
	}
	clientConn, err := grpc.Dial(ms.lis.Addr().String(), ms.dialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial to mock server on %q: %w", ms.lis.Addr(), err)
	}
	for {
		s := clientConn.GetState()
		if s == connectivity.Ready {
			break
		}
		if !clientConn.WaitForStateChange(ctx, s) {
			return nil, nil, ctx.Err()
		}
	}
	return clientConn, testpb.NewTestServiceClient(clientConn), nil
}

var errUnimplemented = grpc.Errorf(codes.Unimplemented, "Unimplemented")

// SetEmptyCallHandler sets the handler for the EmptyCall RPC.
func (ms *MockServer) SetEmptyCallHandler(handler func(context.Context, *testpb.Empty) (*testpb.Empty, error)) {
	ms.emptyCall = handler
}

// EmptyCall calls the handler set with SetEmptyCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) EmptyCall(ctx context.Context, req *testpb.Empty) (*testpb.Empty, error) {
	handler := ms.emptyCall
	if handler == nil {
		return nil, errUnimplemented
	}
	return handler(ctx, req)
}

// SetUnaryCallHandler sets the handler for the UnaryCall RPC.
func (ms *MockServer) SetUnaryCallHandler(handler func(context.Context, *testpb.SimpleRequest) (*testpb.SimpleResponse, error)) {
	ms.unaryCall = handler
}

// UnaryCall calls the handler set with SetUnaryCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) UnaryCall(ctx context.Context, req *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
	handler := ms.unaryCall
	if handler == nil {
		return nil, errUnimplemented
	}
	return handler(ctx, req)
}

// SetStreamingOutputCallHandler sets the handler for the StreamingOutputCall RPC.
func (ms *MockServer) SetStreamingOutputCallHandler(handler func(*testpb.StreamingOutputCallRequest, testpb.TestService_StreamingOutputCallServer) error) {
	ms.streamingOutputCall = handler
}

// StreamingOutputCall calls the handler set with SetStreamingOutputCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) StreamingOutputCall(req *testpb.StreamingOutputCallRequest, stream testpb.TestService_StreamingOutputCallServer) error {
	handler := ms.streamingOutputCall
	if handler == nil {
		return errUnimplemented
	}
	return handler(req, stream)
}

// SetStreamingInputCallHandler sets the handler for the StreamingInputCall RPC.
func (ms *MockServer) SetStreamingInputCallHandler(handler func(testpb.TestService_StreamingInputCallServer) error) {
	ms.streamingInputCall = handler
}

// StreamingInputCall calls the handler set with SetStreamingInputCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) StreamingInputCall(stream testpb.TestService_StreamingInputCallServer) error {
	handler := ms.streamingInputCall
	if handler == nil {
		return errUnimplemented
	}
	return handler(stream)
}

// SetFullDuplexCallHandler sets the handler for the FullDuplexCall RPC.
func (ms *MockServer) SetFullDuplexCallHandler(handler func(testpb.TestService_FullDuplexCallServer) error) {
	ms.fullDuplexCall = handler
}

// FullDuplexCall calls the handler set with SetFullDuplexCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) FullDuplexCall(stream testpb.TestService_FullDuplexCallServer) error {
	handler := ms.fullDuplexCall
	if handler == nil {
		return errUnimplemented
	}
	return handler(stream)
}

// SetHalfDuplexCallHandler sets the handler for the HalfDuplexCall RPC.
func (ms *MockServer) SetHalfDuplexCallHandler(handler func(testpb.TestService_HalfDuplexCallServer) error) {
	ms.halfDuplexCall = handler
}

// HalfDuplexCall calls the handler set with SetHalfDuplexCallHandler.
// If that handler is nil, it delegates the call to the UnimplementedTestServiceServer.
func (ms *MockServer) HalfDuplexCall(stream testpb.TestService_HalfDuplexCallServer) error {
	handler := ms.halfDuplexCall
	if handler == nil {
		return errUnimplemented
	}
	return handler(stream)
}
