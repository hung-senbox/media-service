package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var basePath string

func init() {
	// Lấy thư mục hiện tại của project (ví dụ: /Applications/Senbox/src/term-service)
	wd, _ := os.Getwd()

	// Gốc log = thư mục cha (/Applications/Senbox/src) + /Log/Term-Service
	srcDir := filepath.Dir(wd)
	basePath = filepath.Join(srcDir, "Log")
}

// parseLogLevel chuyển string → logrus.Level
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

// initLogger tạo logger theo thư mục con
func initLogger(subFolder string, level logrus.Level) (*logrus.Logger, error) {
	logDir := filepath.Join(basePath, subFolder)
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, "app.log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetOutput(f)
	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	logger.SetLevel(level)

	return logger, nil
}

// WriteLogMsg → ghi vào LogMsg/app.log
func WriteLogMsg(levelStr string, msg string) {
	logger, err := initLogger("LogMsg", parseLogLevel(levelStr))
	if err != nil {
		fmt.Println("Lỗi khởi tạo logger:", err)
		return
	}
	logger.Log(parseLogLevel(levelStr), msg)
}

// WriteLogData → ghi vào LogData/app.log
func WriteLogData(levelStr string, data any) {
	logger, err := initLogger("LogData", parseLogLevel(levelStr))
	if err != nil {
		fmt.Println("Lỗi khởi tạo logger:", err)
		return
	}
	logger.WithField("data", data).Log(parseLogLevel(levelStr), "Log dữ liệu")
}

// WriteLogEx → ghi vào LogEx/app.log
func WriteLogEx(levelStr string, msg string, data any) {
	logger, err := initLogger("LogEx", parseLogLevel(levelStr))
	if err != nil {
		fmt.Println("Lỗi khởi tạo logger:", err)
		return
	}
	logger.WithField("data", data).Log(parseLogLevel(levelStr), msg)
}
