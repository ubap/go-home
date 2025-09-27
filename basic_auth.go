package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// User represents a simple user with basic auth credentials
type BasicUser struct {
	Username     string
	PasswordHash string
	Role         string // "admin" or "user"
	LastAccess   time.Time
}

// BasicAuthManager handles simple authentication
type BasicAuthManager struct {
	users map[string]*BasicUser
	realm string
}

// NewBasicAuthManager creates a new basic auth manager
func NewBasicAuthManager() *BasicAuthManager {
	auth := &BasicAuthManager{
		users: make(map[string]*BasicUser),
		realm: "Home Automation System",
	}

	// Create default users (change these passwords!)
	auth.AddUser("admin", "admin123!", "admin")
	auth.AddUser("family", "family123!", "user")

	fmt.Println("Basic Auth users created:")
	fmt.Println("Admin: username=admin, password=admin123!")
	fmt.Println("User: username=family, password=family123!")
	fmt.Println("IMPORTANT: Change these passwords immediately!")

	return auth
}

// AddUser adds a new user to the system
func (ba *BasicAuthManager) AddUser(username, password, role string) {
	hash := sha256.Sum256([]byte(password))
	ba.users[username] = &BasicUser{
		Username:     username,
		PasswordHash: base64.StdEncoding.EncodeToString(hash[:]),
		Role:         role,
		LastAccess:   time.Now(),
	}
}

// ValidateCredentials checks if the provided credentials are valid
func (ba *BasicAuthManager) ValidateCredentials(username, password string) (*BasicUser, bool) {
	user, exists := ba.users[username]
	if !exists {
		return nil, false
	}

	// Hash the provided password
	hash := sha256.Sum256([]byte(password))
	providedHash := base64.StdEncoding.EncodeToString(hash[:])

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(user.PasswordHash), []byte(providedHash)) == 1 {
		user.LastAccess = time.Now()
		return user, true
	}

	return nil, false
}

// BasicAuthMiddleware provides HTTP Basic Authentication
func (ba *BasicAuthManager) BasicAuthMiddleware(requiredRole string) func(http.Handler) http.Handler {
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

			// Check role permissions
			if !ba.hasPermission(user.Role, requiredRole) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user info to headers for use in handlers
			r.Header.Set("X-Auth-User", user.Username)
			r.Header.Set("X-Auth-Role", user.Role)

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// parseBasicAuth parses the Basic Auth header
func (ba *BasicAuthManager) parseBasicAuth(auth string) (username, password string, ok bool) {
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

// requestAuth sends a 401 response requesting authentication
func (ba *BasicAuthManager) requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, ba.realm))
	http.Error(w, "Authentication required", http.StatusUnauthorized)
}

// hasPermission checks if a role has the required permission
func (ba *BasicAuthManager) hasPermission(userRole, requiredRole string) bool {
	// Admin can do everything
	if userRole == "admin" {
		return true
	}

	// User can access user-level endpoints
	if requiredRole == "user" && userRole == "user" {
		return true
	}

	return false
}
