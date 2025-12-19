package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// ServerConfig holds web server configuration
type ServerConfig struct {
	Enabled         bool
	Port            string
	Host            string
	JWTSecret       []byte
	JWTExpiry       time.Duration
	CORSOrigins     []string
	CORSCredentials bool
}

// GetServerConfig loads server configuration from environment variables
func GetServerConfig() ServerConfig {
	// Web server enabled flag
	enabled, _ := strconv.ParseBool(os.Getenv("WEB_ENABLED"))
	if os.Getenv("WEB_ENABLED") == "" {
		enabled = false // Default to disabled
	}

	// Web server port
	port := os.Getenv("WEB_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Web server host
	host := os.Getenv("WEB_HOST")
	if host == "" {
		host = "0.0.0.0" // Default to all interfaces
	}

	// JWT secret key
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-this-in-production"
	}

	// JWT expiry duration
	expiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if expiryHours == 0 {
		expiryHours = 24 // Default to 24 hours
	}

	// CORS allowed origins
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var origins []string
	if originsStr != "" {
		// Split by comma and trim spaces
		parts := strings.Split(originsStr, ",")
		for _, origin := range parts {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}
	if len(origins) == 0 {
		origins = []string{"*"} // Default to allow all (not recommended for production)
	}

	// CORS allow credentials
	credentials, _ := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))

	return ServerConfig{
		Enabled:         enabled,
		Port:            port,
		Host:            host,
		JWTSecret:       []byte(secret),
		JWTExpiry:       time.Duration(expiryHours) * time.Hour,
		CORSOrigins:     origins,
		CORSCredentials: credentials,
	}
}
