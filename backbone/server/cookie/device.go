package cookie

import (
	"context"
	"errors"
	"net/http"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	duidAttributeKey = attribute.Key("duid")
	duidContextKey   = contextKey("duid")
)

// SetDeviceUID sets (replaces) the user device UID cookie.
func (m *Middleware) SetDeviceUID(w http.ResponseWriter, value ksuid.KSUID) {
	c := *m.duidCookieSettings
	c.Value = value.String()
	http.SetCookie(w, &c)
}

// DeviceUID is a middleware that reads device UID cookies from requests,
// or generates a new device UID.
func (m *Middleware) DeviceUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		duidCookie, err := r.Cookie(m.duidCookieSettings.Name)
		var uid ksuid.KSUID
		switch {
		case errors.Is(err, http.ErrNoCookie):
			span.AddEvent("set new duid cookie")
			uid = m.MustNewUID()
			m.SetDeviceUID(w, uid)
		case err != nil:
			span.RecordError(err)
			uid = m.MustNewUID()
			m.SetDeviceUID(w, uid)
		default:
			uid, err = m.parseUID(duidCookie.Value)
			if err != nil {
				span.AddEvent("replace invalid duid cookie")
				uid = m.MustNewUID()
				m.SetDeviceUID(w, uid)
			}
		}
		span.SetAttributes(duidAttributeKey.String(uid.String()))
		ctx := context.WithValue(r.Context(), duidContextKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DeviceUIDFromContext returns the device UID from a request context,
// or returns ksuid.Nil if there is no device UID present.
func DeviceUIDFromContext(ctx context.Context) ksuid.KSUID {
	if uid, ok := ctx.Value(duidContextKey).(ksuid.KSUID); ok {
		return uid
	}
	return ksuid.Nil
}
