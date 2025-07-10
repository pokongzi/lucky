package log

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	SetLevel(logrus.DebugLevel)
	SetOutput("./log/shmily.log", false)

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
