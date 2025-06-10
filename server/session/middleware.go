package session

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"
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
		if raw := r.Header.Get("Cookie"); raw != "" {
			cookies, err := http.ParseCookie(raw)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid cookie header: %v", err), http.StatusBadRequest)
				return
			}

			cookies = slices.DeleteFunc(cookies, func(c *http.Cookie) bool {
				return c.Name != AuthCookieName
			})

			if len(cookies) > 1 {
				http.Error(w, "Multiple authentication cookies", http.StatusUnauthorized)
				return
			} else if len(cookies) == 1 {
				rawToken := cookies[0].Value
				var s Session
				if !s.ValidateJWT(rawToken) {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				if s.IsExpired() {
					http.Error(w, "Expired session", http.StatusUnauthorized)
					return
				}

				s.UpdateExpiration()
				http.SetCookie(w, &http.Cookie{
					Name:     AuthCookieName,
					Value:    s.CreateJWT(),
					MaxAge:   int(time.Until(s.Expire).Seconds()),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteDefaultMode,
				})

				ctx := context.WithValue(r.Context(), sessionContextKey{}, s)
				r = r.WithContext(ctx)
			}
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
		raw := r.Header.Get("Cookie")
		if raw == "" {
			http.Error(w, "Unauthorized, no cookie found", http.StatusUnauthorized)
			return
		}

		cookies, err := http.ParseCookie(raw)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid cookie header: %v", err), http.StatusUnauthorized)
			return
		}

		cookies = slices.DeleteFunc(cookies, func(c *http.Cookie) bool {
			return c.Name != AuthCookieName
		})

		if len(cookies) == 0 {
			http.Error(w, "No authentication cookie set", http.StatusUnauthorized)
			return
		} else if len(cookies) > 1 {
			http.Error(w, "Multiple authentication cookies", http.StatusUnauthorized)
			return
		}

		rawToken := cookies[0].Value
		var s Session
		if !s.ValidateJWT(rawToken) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if s.IsExpired() {
			http.Error(w, "Expired session", http.StatusUnauthorized)
			return
		}

		s.UpdateExpiration()
		http.SetCookie(w, &http.Cookie{
			Name:     AuthCookieName,
			Value:    s.CreateJWT(),
			MaxAge:   int(time.Until(s.Expire).Seconds()),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		})

		ctx := context.WithValue(r.Context(), sessionContextKey{}, s)
		r = r.WithContext(ctx)

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
		if raw := r.Header.Get("Cookie"); raw != "" {
			cookies, err := http.ParseCookie(raw)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid cookie header: %v", err), http.StatusBadRequest)
				return
			}

			cookies = slices.DeleteFunc(cookies, func(c *http.Cookie) bool {
				return c.Name != AuthCookieName
			})

			if len(cookies) == 1 {
				http.Error(w, "User must not be authorized", http.StatusBadRequest)
				return
			} else if len(cookies) > 1 {
				http.Error(w, "Multiple authentication cookies", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), sessionContextKey{}, Session{})
			r = r.WithContext(ctx)
		}

		h.ServeHTTP(w, r)
	})
}

// See [NoAuthMiddleware]
func NoAuthMiddlewareFunc(f http.HandlerFunc) http.Handler {
	return NoAuthMiddleware(f)
}
