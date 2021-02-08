package util

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var once sync.Once

// 单例模式
func GetLogger() *logrus.Logger {
	once.Do(func() {
		logger = newLogger()
	})
	return logger
}
func newLogger() *logrus.Logger {
	// 获取日志写入文件
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println("Fail to create logs:", err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)

	// 创建日志文件
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println("Fail to create log file:", err.Error())
		}
	}

	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Fail to write logs:", err.Error())
	}

	// 实例化logger
	logger := logrus.New()
	logger.Out = src
	fmt.Println("setting log level:", LogLevel)
	switch LogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
		fmt.Println("No available logger level! Set default: debug level")
	}
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: LogFormat,
	})
	return logger
}
