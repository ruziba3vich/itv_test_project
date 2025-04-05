package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is a custom wrapper around logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger initializes a new logger with file output
func NewLogger() (*Logger, error) {
	log := logrus.New()

	// Set output to app.log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	// Set log format (JSON for structured logging, optional)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set log level (can be configured later via env vars)
	log.SetLevel(logrus.InfoLevel)

	return &Logger{log}, nil
}

// Convenience methods for leveled logging
func (l *Logger) Debug(msg string, fields ...map[string]any) {
	l.addFields(l.Logger.Debug, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Info(msg)
	} else {
		l.Logger.Info(msg)
	}
}

func (l *Logger) Warn(msg string, fields ...map[string]any) {
	l.addFields(l.Logger.Warn, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...map[string]any) {
	l.addFields(l.Logger.Error, msg, fields...)
}

// addFields adds optional key-value pairs to the log entry
func (l *Logger) addFields(logFunc func(...any), msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		entry := l.WithFields(fields[0])
		entry.Log(logrus.InfoLevel, msg) // Still not dynamic
	} else {
		logFunc(msg)
	}
}
