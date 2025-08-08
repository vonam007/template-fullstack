package middleware

import (
	"time"

	"template-fullstack/backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Str("user_agent", param.Request.UserAgent()).
			Msg("HTTP Request")
		return ""
	})
}

func StructuredLoggingMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)

		logEvent := log.Info()
		if c.Writer.Status() >= 400 {
			logEvent = log.Error()
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size()).
			Msg("HTTP Request")
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			log := logger.New()
			log.Error().
				Err(err.Err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("Request processing error")

			// If no response has been written yet, write error response
			if !c.Writer.Written() {
				c.JSON(500, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error",
					},
				})
			}
		}
	}
}
