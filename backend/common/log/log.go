package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"lucky/common/config"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	initFromConfig()
}

// initFromConfig 从配置文件初始化日志设置
func initFromConfig() {
	// 读取日志配置
	logSection := config.Config.Section("log")

	// 设置日志等级
	levelStr := logSection.Key("level").String()
	level := parseLogLevel(levelStr)
	SetLevel(level)

	// 设置日志文件路径
	logFile := logSection.Key("file").String()
	if logFile == "" {
		logFile = "./log/lucky.log" // 默认路径
	}

	// 确保日志目录存在
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	// 设置控制台输出
	withConsoleStr := logSection.Key("with_console").String()
	withConsole := strings.ToLower(withConsoleStr) == "true"

	SetOutput(logFile, withConsole)

	// 设置日志格式
	SetFormatter()
}

// parseLogLevel 解析日志等级字符串
func parseLogLevel(levelStr string) logrus.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel // 默认debug级别
	}
}

// SetFormatter 设置日志格式
func SetFormatter() {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	}
	logger.SetFormatter(formatter)
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func SetOutput(logFile string, withConsole bool) {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + logFile)
	}
	if withConsole {
		logger.SetOutput(io.MultiWriter(f, os.Stdout))
	} else {
		logger.SetOutput(f)
	}
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func AddHook(hook logrus.Hook) {
	logger.AddHook(hook)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}
