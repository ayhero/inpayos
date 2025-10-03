package log

import (
	"fmt"
	"inpayos/internal/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Service-specific loggers - using sync.Map for better concurrent performance
var serviceLoggers sync.Map

const (
	DefaultLogger = ""
)

func Get() *logrus.Logger {
	return GetServiceLogger(DefaultLogger)
}

func GetServiceLogger(service string) *logrus.Logger {
	today := time.Now().Format(time.DateOnly)
	loggerKey := service
	if service != "" {
		loggerKey = fmt.Sprintf("%s-", loggerKey)
	}
	loggerKey = fmt.Sprintf("%s%s", loggerKey, today)
	return GetServiceLoggerWithoutDate(loggerKey)
}

func GetServiceLoggerWithoutDate(service string) *logrus.Logger {
	loggerKey := service

	// Try to load existing logger
	if value, exists := serviceLoggers.Load(loggerKey); exists {
		return value.(*logrus.Logger)
	}

	// Create new logger
	cfg := config.Get().Log
	logger := initLogger(cfg.Path, cfg.Level, service)

	// Store or load existing (handles race condition automatically)
	if actual, loaded := serviceLoggers.LoadOrStore(loggerKey, logger); loaded {
		// Another goroutine created it first, use that one
		return actual.(*logrus.Logger)
	}

	// We created it successfully
	return logger
}

func Init() {
	cfg := config.Get().Log
	if cfg == nil {
		return
	}
	serviceLoggers.Store(DefaultLogger, initLogger(cfg.Path, cfg.Level, DefaultLogger))
	for name, scfg := range cfg.Services {
		serviceLoggers.Store(name, initLogger(scfg.Path, scfg.Level, name))
	}
}

func initLogger(path string, level string, service string) (logger *logrus.Logger) {
	logger = logrus.New()
	// 获取当前日期
	t := time.Now()
	dateStr := t.Format(time.DateOnly)

	// 如果没有配置日志路径，使用默认路径
	if path == "" {
		// 使用相对路径，确保在当前工作目录下创建日志
		path = "logs"
	}

	// 确保日志路径是绝对路径
	if !filepath.IsAbs(path) {
		currentDir, err := os.Getwd()
		if err == nil {
			path = filepath.Join(currentDir, path)
		}
	}

	// 如果日志路径不存在，则创建
	if err := os.MkdirAll(path, 0o755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		// 降级到当前目录
		path = "logs"
		os.MkdirAll(path, 0o755)
	}

	// 日志文件名格式：service-date.log
	var logFileName string
	if service == "" {
		logFileName = filepath.Join(path, fmt.Sprintf("%s.log", dateStr))
	} else {
		logFileName = filepath.Join(path, fmt.Sprintf("%s-%s.log", service, dateStr))
	}

	// 打开或创建日志文件
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		logger.SetOutput(os.Stdout)
	} else {
		// 设置多重输出，同时写入文件和标准输出
		mw := io.MultiWriter(os.Stdout, logFile)
		logger.SetOutput(mw)
	}

	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 设置日志级别
	if level == "" {
		level = "debug"
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(lvl)
	}

	return
}
