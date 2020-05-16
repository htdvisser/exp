package grpctest

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	testpb "google.golang.org/grpc/test/grpc_testing"
)

func TestMockServer(t *testing.T) {
	mock := NewMockServer(nil, nil)

	_, _, err := mock.ConnectTestClient(context.Background())
	if err == nil {
		t.Error("mock.ConnectTestClient(ctx) did not return error while not started yet")
	}

	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mock.Start(startCtx); err != nil {
		t.Fatal(err)
	}
	defer mock.Stop()

	if err := mock.Start(startCtx); err == nil {
		t.Error("mock.Start(ctx) returned did not return error while already started")
	}

	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cc, client, err := mock.ConnectTestClient(connectCtx)
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	t.Run("Unimplemented", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err = client.EmptyCall(callCtx, &testpb.Empty{})
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.EmptyCall returned %q, expected Unimplemented error", err)
			}
		})

		t.Run("Unary", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err = client.UnaryCall(callCtx, &testpb.SimpleRequest{})
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.UnaryCall returned %q, expected Unimplemented error", err)
			}
		})

		t.Run("StreamingOutput", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.StreamingOutputCall(callCtx, &testpb.StreamingOutputCallRequest{})
			if err != nil {
				t.Errorf("client.StreamingOutputCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.StreamingOutputCall returned %q, expected Unimplemented error", err)
			}
		})

		t.Run("StreamingInput", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.StreamingInputCall(callCtx)
			if err != nil {
				t.Errorf("client.StreamingInputCall failed too early: %v", err)
			}
			_, err = stream.CloseAndRecv()
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.StreamingInputCall returned %q, expected Unimplemented error", err)
			}
		})

		t.Run("FullDuplex", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.FullDuplexCall(callCtx)
			if err != nil {
				t.Errorf("client.FullDuplexCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.FullDuplexCall returned %q, expected Unimplemented error", err)
			}
		})

		t.Run("HalfDuplex", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.HalfDuplexCall(callCtx)
			if err != nil {
				t.Errorf("client.HalfDuplexCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unimplemented {
				t.Errorf("client.HalfDuplexCall returned %q, expected Unimplemented error", err)
			}
		})
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		mock.SetEmptyCallHandler(func(context.Context, *testpb.Empty) (*testpb.Empty, error) {
			return nil, status.Errorf(codes.Unauthenticated, "EmptyCall requires authentication")
		})
		mock.SetUnaryCallHandler(func(context.Context, *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
			return nil, status.Errorf(codes.Unauthenticated, "UnaryCall requires authentication")
		})
		mock.SetStreamingOutputCallHandler(func(*testpb.StreamingOutputCallRequest, testpb.TestService_StreamingOutputCallServer) error {
			return status.Errorf(codes.Unauthenticated, "StreamingOutputCall requires authentication")
		})
		mock.SetStreamingInputCallHandler(func(testpb.TestService_StreamingInputCallServer) error {
			return status.Errorf(codes.Unauthenticated, "StreamingInputCall requires authentication")
		})
		mock.SetFullDuplexCallHandler(func(testpb.TestService_FullDuplexCallServer) error {
			return status.Errorf(codes.Unauthenticated, "FullDuplexCall requires authentication")
		})
		mock.SetHalfDuplexCallHandler(func(testpb.TestService_HalfDuplexCallServer) error {
			return status.Errorf(codes.Unauthenticated, "HalfDuplexCall requires authentication")
		})

		t.Run("Empty", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err = client.EmptyCall(callCtx, &testpb.Empty{})
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.EmptyCall returned %q, expected Unauthenticated error", err)
			}
		})

		t.Run("Unary", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err = client.UnaryCall(callCtx, &testpb.SimpleRequest{})
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.UnaryCall returned %q, expected Unauthenticated error", err)
			}
		})

		t.Run("StreamingOutput", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.StreamingOutputCall(callCtx, &testpb.StreamingOutputCallRequest{})
			if err != nil {
				t.Errorf("client.StreamingOutputCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.StreamingOutputCall returned %q, expected Unauthenticated error", err)
			}
		})

		t.Run("StreamingInput", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.StreamingInputCall(callCtx)
			if err != nil {
				t.Errorf("client.StreamingInputCall failed too early: %v", err)
			}
			_, err = stream.CloseAndRecv()
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.StreamingInputCall returned %q, expected Unauthenticated error", err)
			}
		})

		t.Run("FullDuplex", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.FullDuplexCall(callCtx)
			if err != nil {
				t.Errorf("client.FullDuplexCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.FullDuplexCall returned %q, expected Unauthenticated error", err)
			}
		})

		t.Run("HalfDuplex", func(t *testing.T) {
			callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.HalfDuplexCall(callCtx)
			if err != nil {
				t.Errorf("client.HalfDuplexCall failed too early: %v", err)
			}
			_, err = stream.Recv()
			if status.Code(err) != codes.Unauthenticated {
				t.Errorf("client.HalfDuplexCall returned %q, expected Unauthenticated error", err)
			}
		})
	})
}
