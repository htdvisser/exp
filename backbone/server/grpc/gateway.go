package grpc

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func handleProtoError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	if err == runtime.ErrUnknownURI {
		http.NotFound(w, r)
		return
	}
	runtime.DefaultHTTPProtoErrorHandler(ctx, mux, marshaler, w, r, err)
}

func handleStreamError(ctx context.Context, err error) *runtime.StreamError {
	return runtime.DefaultHTTPStreamErrorHandler(ctx, err)
}
