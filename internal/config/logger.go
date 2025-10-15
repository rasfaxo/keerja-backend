package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// InitLogger initializes the global logger instance
func InitLogger() {
	logger = logrus.New()

	// Set log format based on environment
	env := os.Getenv("APP_ENV")
	if env == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		logger.SetLevel(logrus.DebugLevel)
	}

	// Set output to stdout
	logger.SetOutput(os.Stdout)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

// WithField adds a single field to the logger
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields adds multiple fields to the logger
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithError adds an error field to the logger
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}
