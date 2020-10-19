package rjs

import (
	"context"
	"fmt"
)

// Client is the interface that makes requests to the RJS API.
//
// uri is the root-relative URI of the API call. Example: "/api/v1/router/info".
// request is the request message. The implementation sends this object in the request body.
// response is the response message. The implementation decodes the response body into this object.
// Returned errors can wrap errors from package encoding/json if JSON encoding/decoding failed,
// *url.Error if the request failed, or *Error if the server did not accept the request.
type Client interface {
	Call(ctx context.Context, uri string, request interface{}, response interface{}) error
}

// Error may be returned (possibly wrapped) by implementations of the Client interface.
type Error struct {
	StatusCode int
	Body       []byte
}

func (err *Error) Error() string {
	return fmt.Sprintf("RJS server returned %d", err.StatusCode)
}
