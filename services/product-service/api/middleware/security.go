// api/middleware/security.go
package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"

	"github.com/ntdt/product-service/config"
	"github.com/ntdt/product-service/pkg/logger"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = NewUUID()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// NewUUID generates a unique UUID string
func NewUUID() string {
	return "req-" + time.Now().Format("20060102150405") + "-" + RandString(8)
}

// RandString generates a random string of specified length
func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}

// Logger middleware for logging requests
func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		requestID, _ := c.Get("RequestID")

		log.Info("Request processed",
			logger.Fields{
				"status":     statusCode,
				"latency":    latency,
				"client_ip":  clientIP,
				"method":     method,
				"path":       path,
				"request_id": requestID,
			},
		)
	}
}

// Cors middleware for handling CORS
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter middleware for rate limiting requests
func RateLimiter(cfg config.RateLimitConfig) gin.HandlerFunc {
	// Create a new rate limiter for each IP address
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Create a new limiter if it doesn't exist
		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(rate.Limit(cfg.Requests), cfg.Requests)
		}

		// Check if the request exceeds the rate limit
		if !limiters[ip].Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": http.StatusTooManyRequests,
				"error":  "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Auth middleware for JWT authentication
func Auth(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header format is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Invalid authorization format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user ID in context
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("UserID", userID)
		}

		c.Next()
	}
}
