package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

// NewLogger creates a new instance of the Logger
func New() *Logger {

	logger := &Logger{}
	logger.init()
	return logger
}

func (s *Logger) Log() *zap.Logger {
	return s.log
}

func (logger *Logger) init() {
	// Initialize a new Zap logger instance
	var err error
	logger.log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.log.Sync() // Flushes buffer, if any
}
