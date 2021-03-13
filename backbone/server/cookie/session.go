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
	suidAttributeKey = attribute.Key("suid")
	suidContextKey   = contextKey("suid")
)

func (m *Middleware) setSessionUID(w http.ResponseWriter, value ksuid.KSUID) {
	c := *m.suidCookieSettings
	c.Value = value.String()
	http.SetCookie(w, &c)
}

// SessionUID is a middleware that reads session UID cookies from requests,
// or generates a new session UID.
func (m *Middleware) SessionUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		suidCookie, err := r.Cookie(m.suidCookieSettings.Name)
		var uid ksuid.KSUID
		switch {
		case errors.Is(err, http.ErrNoCookie):
			span.AddEvent("set new suid cookie")
			uid = m.MustNewUID()
			m.setSessionUID(w, uid)
		case err != nil:
			span.RecordError(err)
			uid = m.MustNewUID()
			m.setSessionUID(w, uid)
		default:
			uid, err = m.parseUID(suidCookie.Value)
			if err != nil {
				span.AddEvent("replace invalid suid cookie")
				uid = m.MustNewUID()
				m.setSessionUID(w, uid)
			}
		}
		span.SetAttributes(suidAttributeKey.String(uid.String()))
		ctx := context.WithValue(r.Context(), suidContextKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SessionUIDFromContext returns the session UID from a request context,
// or returns ksuid.Nil if there is no session UID present.
func SessionUIDFromContext(ctx context.Context) ksuid.KSUID {
	if uid, ok := ctx.Value(suidContextKey).(ksuid.KSUID); ok {
		return uid
	}
	return ksuid.Nil
}
