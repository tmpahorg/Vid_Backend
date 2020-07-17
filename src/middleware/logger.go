package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Now()
		latency := float64(stop.Sub(start).Nanoseconds())

		method := c.Request.Method
		path := c.Request.URL.Path
		code := c.Writer.Status()
		length := math.Abs(float64(c.Writer.Size()))
		ip := c.ClientIP()

		entry := logger.WithFields(logrus.Fields{
			"Module":   "gin",
			"Method":   method,
			"Path":     path,
			"Code":     code,
			"Length":   length,
			"ClientIP": ip,
			"Latency":  latency,
		})
		if len(c.Errors) == 0 {
			msg := fmt.Sprintf("[Gin] %3d | %12s | %15s | %8s | %-7s %s",
				code, xnumber.RenderLatency(latency), ip, xnumber.RenderByte(length), method, path)
			if code >= 500 {
				entry.Error(msg)
			} else if code >= 400 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		} else {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}
