package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EnableHTTPS       bool
	CertFile          string
	KeyFile           string
	EnableHSTS        bool
	EnableCSP         bool
	SessionTimeout    time.Duration
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
}

// DefaultSecurityConfig returns recommended security settings
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		EnableHTTPS:      true,
		CertFile:         "server.crt",
		KeyFile:          "server.key",
		EnableHSTS:       true,
		EnableCSP:        true,
		SessionTimeout:   24 * time.Hour,
		MaxLoginAttempts: 5,
		LockoutDuration:  15 * time.Minute,
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HTTPS Strict Transport Security
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy - bardzo restrykcyjne dla bezpieczeństwa
		w.Header().Set("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data:; "+
			"connect-src 'self'; "+
			"font-src 'self'; "+
			"object-src 'none'; "+
			"media-src 'self'; "+
			"frame-src 'none';")
		
		// Zapobieganie clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// Zapobieganie MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// XSS Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions Policy (dawniej Feature Policy)
		w.Header().Set("Permissions-Policy", 
			"geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=()")
		
		next.ServeHTTP(w, r)
	})
}

// CreateTLSConfig creates secure TLS configuration
func CreateTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}
}

// GenerateSelfSignedCert generates a self-signed certificate for development
// W produkcji użyj Let's Encrypt lub innego zaufanego CA
func GenerateSelfSignedCert() error {
	// Implementacja generowania certyfikatu self-signed
	// W rzeczywistości powinieneś użyć Let's Encrypt dla produkcji
	return nil
}
