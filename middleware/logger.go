package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

func Loggoer() gin.HandlerFunc {
	filepath := "log/log"
	linkName := "latest_log.log"
	src, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}
	logger := logrus.New()
	logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//日志分割
	logWriter, _ := retalog.New(
		filepath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
		retalog.WithLinkName(linkName),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(Hook)
	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{}))
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Milliseconds())/1000000.0)))
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unKnow"
		}
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize > 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.URL.Path

		entry := logger.WithFields(logrus.Fields{
			"HostName":   hostName,
			"ClientIP":   clientIP,
			"UserAgent":  userAgent,
			"Method":     method,
			"Path":       path,
			"SPendTime":  spendTime,
			"DataSize":   dataSize,
			"StatusCode": statusCode,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
