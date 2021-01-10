package grpc

import (
	"context"
	"errors"
	"net/http"
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

var defaultRequestHeaders = []string{
	"accept-language",
	"authorization",
	"cookie",
	"correlation-context",
	"referer",
	"traceparent",
	"user-agent",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-spanid",
	"x-b3-traceid",
	"x-forwarded-proto",
	"x-request-id",
}

type runtimeHeaders map[string]string

func (h runtimeHeaders) add(headers ...string) {
	for _, header := range headers {
		header = textproto.CanonicalMIMEHeaderKey(header)
		h[header] = header
	}
}

func (h runtimeHeaders) match(header string) (string, bool) {
	out, ok := h[header]
	return out, ok
}

func handleError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, runtime.ErrNotMatch) {
		http.NotFound(w, r)
		return
	}
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func handleStreamError(ctx context.Context, err error) *status.Status {
	return runtime.DefaultStreamErrorHandler(ctx, err)
}
