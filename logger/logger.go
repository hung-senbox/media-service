package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers   = make(map[string]*logrus.Logger)
	once      sync.Once
	basePath  string
	initMutex sync.Mutex
)

// parse log level
func parseLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// setup base path only once
func initBasePath() {
	once.Do(func() {
		// log path mới
		basePath = filepath.Join("logs", "media-service")

		// tạo thư mục nếu chưa có
		_ = os.MkdirAll(basePath, 0755)
	})
}

func getLogger(folder string) *logrus.Logger {
	initBasePath()

	initMutex.Lock()
	defer initMutex.Unlock()

	if logger, exists := loggers[folder]; exists {
		return logger
	}

	logDir := filepath.Join(basePath, folder)
	os.MkdirAll(logDir, 0755)

	// file rotation config
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	// Also log to console (best for containers)
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	loggers[folder] = logger
	return logger
}

// ------------------------------------------------------
// Public usage
// ------------------------------------------------------

func WriteLogMsg(level string, msg string) {
	getLogger("LogMsg").Log(parseLogLevel(level), msg)
}

func WriteLogData(level string, data any) {
	getLogger("LogData").WithField("data", data).Log(parseLogLevel(level), "Data log")
}

func WriteLogEx(level string, msg string, data any) {
	getLogger("LogEx").WithField("data", data).Log(parseLogLevel(level), msg)
}
