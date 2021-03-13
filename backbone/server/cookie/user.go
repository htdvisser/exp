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
	usuidAttributeKey = attribute.Key("usuid")
	usuidContextKey   = contextKey("usuid")
)

// SetUserSessionUID sets the user session UID cookie, with expiry after maxAge.
func (m *Middleware) SetUserSessionUID(w http.ResponseWriter, value ksuid.KSUID, maxAge int) {
	c := *m.usuidCookieSettings
	c.Value = value.String()
	c.MaxAge = maxAge
	http.SetCookie(w, &c)
}

// UnsetUserSessionUID unsets the user session UID cookie.
func (m *Middleware) UnsetUserSessionUID(w http.ResponseWriter) {
	c := *m.usuidCookieSettings
	c.MaxAge = -1
	http.SetCookie(w, &c)
}

// UserSessionUID is a middleware that reads the user session UID cookie from requests.
func (m *Middleware) UserSessionUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		usuidCookie, err := r.Cookie(m.usuidCookieSettings.Name)
		var uid ksuid.KSUID
		switch {
		case errors.Is(err, http.ErrNoCookie):
			next.ServeHTTP(w, r)
			return
		case err != nil:
			span.RecordError(err)
		default:
			uid, err = m.parseUID(usuidCookie.Value)
			if err != nil {
				span.AddEvent("ignore invalid usuid cookie")
				uid = ksuid.Nil
			}
		}
		span.SetAttributes(usuidAttributeKey.String(uid.String()))
		ctx := context.WithValue(r.Context(), usuidContextKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserSessionUIDFromContext returns the user session UID from a request context,
// or returns ksuid.Nil if there is no user session UID present.
func UserSessionUIDFromContext(ctx context.Context) ksuid.KSUID {
	if uid, ok := ctx.Value(usuidContextKey).(ksuid.KSUID); ok {
		return uid
	}
	return ksuid.Nil
}
