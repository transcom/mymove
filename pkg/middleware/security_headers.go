package middleware

import (
	"net/http"
)

var securityHeaders = map[string]string{
	// set the HSTS header using values recommended by OWASP
	// https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet#Examples
	"strict-transport-security": "max-age=31536000; includeSubdomains; preload",
	// Sets headers to prevent rendering our page in an iframe, prevents clickjacking
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	"X-Frame-Options": "deny",
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/frame-ancestors
	"Content-Security-Policy": "frame-ancestors 'none'",
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
	"X-XSS-Protection": "1; mode=block",
}

// SecurityHeaders adds a set of standard security headers.
func SecurityHeaders(logger Logger) func(inner http.Handler) http.Handler {
	logger.Debug("SecurityHeaders Middleware used")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range securityHeaders {
				w.Header().Set(k, v)
			}
			inner.ServeHTTP(w, r)
		})
	}
}
