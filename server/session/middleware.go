package session

import (
	"context"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
)

const AuthCookieName = "auth"

type sessionContextKey struct{}

// Get retrieves value of session from ctx.
//
// If no such value exists, returns a zero value [Session], which will satisfy [Session.IsZero]
func Get(ctx context.Context) Session {
	v := ctx.Value(sessionContextKey{})
	if v == nil {
		return Session{}
	}
	return v.(Session)
}

// Middleware wraps h, making it parse and validate sessions.
//
// If request is properly authenticated, [Get](r.Context()) will return a non-zero [Session].
// Otherwise, session will be unset.
//
// Prefer [AuthMiddleware] or [NoAuthMiddleware] if a resource is only available to authorized or unauthorized users.
// [Middleware] should only be used when handler has different behavior depending on authentication status.
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookies := r.CookiesNamed(AuthCookieName); len(cookies) > 0 {
			if len(cookies) > 1 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p M-ware FAIL; err= %s", r, "too many auth cookies")
				return
			}
			rawToken := cookies[0].Value
			var ses Session
			if !ses.ValidateJWT(rawToken) {
				Logout(w)
				logger.Log.Debug("req=%p M-ware LOGOUT; reason= %s", r, "session is invalid JWT")
			} else if ses.IsExpired() {
				Logout(w)
				logger.Log.Debug("req=%p M-ware LOGOUT; reason= %s", r, "session is expired")
			} else {
				logger.Log.Debug("req=%p M-ware OK; updated session", r)
			}
			ctx := context.WithValue(r.Context(), sessionContextKey{}, ses)
			r = r.WithContext(ctx)
		} else {
			logger.Log.Debug("req=%p N-ware OK; no cookie", r)
		}
		h.ServeHTTP(w, r)
	})
}

// See [Middleware]
func MiddlewareFunc(f http.HandlerFunc) http.Handler {
	return Middleware(f)
}

// AuthMiddleware wraps h making it reject all requests which are not authorized.
//
// AuthMiddleware modifies the context of request, [Session] can be retrieved with [Get].
func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookies := r.CookiesNamed(AuthCookieName); len(cookies) > 0 {
			if len(cookies) > 1 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p A-ware FAIL; err= %s", r, "too many auth cookies")
				return
			}
			rawToken := cookies[0].Value
			var ses Session
			if !ses.ValidateJWT(rawToken) {
				Logout(w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p A-ware LOGOUT; reason= %s", r, "session is invalid JWT")
				return
			} else if ses.IsExpired() {
				Logout(w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p A-ware LOGOUT; reason= %s", r, "session is expired")
				return
			} else {
				logger.Log.Debug("req=%p A-ware OK; updated session", r)
			}
			ctx := context.WithValue(r.Context(), sessionContextKey{}, ses)
			r = r.WithContext(ctx)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			logger.Log.Debug("req=%p A-ware FAIL; err= %s", r, "user is not authorized")
			return
		}
		h.ServeHTTP(w, r)
	})
}

// See [AuthMiddleware]
func AuthMiddlewareFunc(f http.HandlerFunc) http.Handler {
	return AuthMiddleware(f)
}

// NoAuthMiddleware wraps h making it reject all requests which have authentication headers set.
//
// NoAuthMiddleware deletes session from the context.
//
// Note that requests with invalid Cookie headers will still be rejected.
func NoAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookies := r.CookiesNamed(AuthCookieName); len(cookies) > 0 {
			if len(cookies) > 1 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p N-ware FAIL; err= %s", r, "too many auth cookies")
				return
			}
			rawToken := cookies[0].Value
			var ses Session
			if !ses.ValidateJWT(rawToken) {
				Logout(w)
				logger.Log.Debug("req=%p N-ware LOGOUT; reason= %s", r, "session is invalid JWT")
			} else if ses.IsExpired() {
				Logout(w)
				logger.Log.Debug("req=%p N-ware LOGOUT; reason= %s", r, "session is expired")
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				logger.Log.Debug("req=%p N-ware FAIL; err= %s", r, "user is authorized")
				return
			}
			ctx := context.WithValue(r.Context(), sessionContextKey{}, ses)
			r = r.WithContext(ctx)
		} else {
			logger.Log.Debug("req=%p N-ware OK; no cookie", r)
		}
		h.ServeHTTP(w, r)
	})
}

// See [NoAuthMiddleware]
func NoAuthMiddlewareFunc(f http.HandlerFunc) http.Handler {
	return NoAuthMiddleware(f)
}
