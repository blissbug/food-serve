package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		//method, path, status, duration
		method := c.Request.Method
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		duration := time.Since(t)
		zap.L().Info("Request", zap.String("method", method), zap.String("path", path), zap.Int("status", status), zap.Duration("duration", duration))
	}
}
