package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth is the middleware for Basic Authentication
func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/people") {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authParts := strings.SplitN(auth, " ", 2)
			if len(authParts) != 2 || authParts[0] != "Basic" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			payload, err := base64.StdEncoding.DecodeString(authParts[1])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 || pair[0] != username || pair[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
