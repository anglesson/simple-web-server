package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anglesson/simple-web-server/internal/config"
)

// SecurityHeaders middleware adds security headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy - mais restritivo
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' https://cdnjs.cloudflare.com; img-src 'self' data:; font-src 'self' https://cdnjs.cloudflare.com; connect-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self';")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// XSS Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// HSTS (only in production)
		if config.AppConfig.IsProduction() {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Remove server information
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")
		
		// Additional security headers
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("X-DNS-Prefetch-Control", "off")

		next.ServeHTTP(w, r)
	})
}

// RateLimiter implements basic rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimitMiddleware applies rate limiting to requests
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		clientIP := getClientIP(r)

		// Check rate limit
		if !rl.isAllowed(clientIP) {
			// Check if it's an API request or HTML request
			acceptHeader := r.Header.Get("Accept")
			if strings.Contains(acceptHeader, "application/json") || strings.Contains(r.URL.Path, "/api/") {
				// API request - return JSON
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "Rate limit exceeded. Please try again later."}`))
			} else {
				// HTML request - redirect with error message based on the route
				var redirectPath string
				if strings.Contains(r.URL.Path, "/forget-password") || strings.Contains(r.URL.Path, "/reset-password") {
					redirectPath = "/forget-password?error=rate_limit_exceeded"
				} else {
					redirectPath = "/login?error=rate_limit_exceeded"
				}
				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isAllowed checks if the request is allowed based on rate limiting rules
func (rl *RateLimiter) isAllowed(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this IP
	requests, exists := rl.requests[clientIP]
	if !exists {
		requests = []time.Time{}
	}

	// Remove old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if we're under the limit
	if len(validRequests) >= rl.limit {
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests

	return true
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check for forwarded headers (in order of preference)
	headers := []string{"X-Forwarded-For", "X-Real-IP", "X-Client-IP", "CF-Connecting-IP"}
	
	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			// X-Forwarded-For can contain multiple IPs, take the first one
			if header == "X-Forwarded-For" {
				if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
					ip = strings.TrimSpace(ip[:commaIndex])
				}
			}
			
			// Validate IP format
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// Fallback to remote address
	// Remove port if present
	if colonIndex := strings.LastIndex(r.RemoteAddr, ":"); colonIndex != -1 {
		ip := r.RemoteAddr[:colonIndex]
		if isValidIP(ip) {
			return ip
		}
	}
	
	// Default fallback
	return "unknown"
}

// isValidIP validates if the string is a valid IP address
func isValidIP(ip string) bool {
	// Basic validation - check if it's not empty and contains dots or colons
	if ip == "" || ip == "unknown" {
		return false
	}
	
	// Check for IPv4 format (contains dots)
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) != 4 {
			return false
		}
		for _, part := range parts {
			if part == "" {
				return false
			}
		}
		return true
	}
	
	// Check for IPv6 format (contains colons)
	if strings.Contains(ip, ":") {
		return true
	}
	
	return false
}

// CleanupRateLimiter periodically cleans up old rate limiting data
func (rl *RateLimiter) CleanupRateLimiter() {
	ticker := time.NewTicker(rl.window)
	go func() {
		for range ticker.C {
			rl.cleanup()
		}
	}()
}

// cleanup removes old rate limiting data
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for ip, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validRequests
		}
	}
}
