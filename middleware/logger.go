package middleware

import (
	"Blog/util"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = util.GetLogger()

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		duration := fmt.Sprintf("%d ms", stop.Milliseconds())

		hostName, err := os.Hostname()
		if err != nil {
			hostName = "Unknown"
		}

		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		method := c.Request.Method
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"Hostname": hostName,
			"Status":   statusCode,
			"Time":     duration,
			"IP":       clientIp,
			"Path":     path,
			"Method":   method,
			"Agent":    userAgent,
			"DataSize": dataSize,
		})

		// 系统内部错误
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		// 状态码错误
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
