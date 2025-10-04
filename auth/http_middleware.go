package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (ba *UserManager) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			auth := r.Header.Get("Authorization")
			if auth == "" {
				ba.requestAuth(w)
				return
			}

			// Parse Basic Auth
			username, password, ok := ba.parseBasicAuth(auth)
			if !ok {
				ba.requestAuth(w)
				return
			}

			// Validate credentials
			user, valid := ba.ValidateCredentials(username, password)
			if !valid {
				ba.requestAuth(w)
				return
			}

			// Add user info to headers for use in handlers
			r.Header.Set("X-Auth-User", user.Username)

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func (ba *UserManager) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return "", "", false
	}

	// Decode base64
	encoded := auth[len(prefix):]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}

	// Split username:password
	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	return parts[0], parts[1], true
}

func (ba *UserManager) requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, ba.realm))
	http.Error(w, "Authentication required", http.StatusUnauthorized)
}
