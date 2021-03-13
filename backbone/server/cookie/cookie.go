// Package cookie provides middleware for working with device, session and user session cookies.
package cookie

import (
	"fmt"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/segmentio/ksuid"
	"htdvisser.dev/exp/backbone/server"
)

// Option is an option for the cookie middleware.
type Option interface {
	applyTo(*Middleware)
}

type optionFunc func(*Middleware)

func (f optionFunc) applyTo(opts *Middleware) {
	f.applyTo(opts)
}

// NewMiddleware returns new cookie middleware.
//
// Options can be used to override the default options that work with:
// a "suid" (session) cookie that expires when the browser is closed,
// a "duid" (device) cookie that expires after 180 days,
// a "usuid" (user session) cookie.
//
// Note that the cookie middleware doesn't prevent (malicious) clients from
// sending cookies after they have expired.
func NewMiddleware(opts ...Option) *Middleware {
	m := &Middleware{
		clock: clock.New(),
		suidCookieSettings: &http.Cookie{
			Name:     "suid",
			Secure:   false, // allow on HTTP.
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
		duidCookieSettings: &http.Cookie{
			Name:     "duid",
			MaxAge:   60 * 60 * 24 * 180,
			Secure:   false, // allow on HTTP.
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
		usuidCookieSettings: &http.Cookie{
			Name:     "usuid",
			Secure:   true, // only allow on HTTPS.
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.applyTo(m)
	}
	return m
}

// Middleware is cookie middleware. Use NewMiddleware to create a new Middleware.
type Middleware struct {
	clock               clock.Clock
	suidCookieSettings  *http.Cookie
	duidCookieSettings  *http.Cookie
	usuidCookieSettings *http.Cookie
}

// MustNewUID returns a new KSUID or panics if it's unable to.
func (m *Middleware) MustNewUID() ksuid.KSUID {
	uid, err := ksuid.NewRandomWithTime(m.clock.Now())
	if err != nil {
		panic(fmt.Errorf("failed to generate KSUID: %w", err))
	}
	return uid
}

func (*Middleware) parseUID(uid string) (ksuid.KSUID, error) {
	return ksuid.Parse(uid)
}

type contextKey string

// Register registers the middleware to a backbone server.
func (m *Middleware) Register(s *server.Server) {
	s.HTTP.AddMiddleware(m.DeviceUID, m.SessionUID, m.UserSessionUID)
	s.InternalHTTP.AddMiddleware(m.DeviceUID, m.SessionUID, m.UserSessionUID)
}
