package omnilogger

import (
	"fmt"
	"omnilogger/config"
	"omnilogger/model"
	pkg "omnilogger/pkg"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	DEBUG config.LogLevel = "DEBUG" // Debug level logging.
	INFO  config.LogLevel = "INFO"  // Information level logging.
	WARN  config.LogLevel = "WARN"  // Warning level logging.
	ERROR config.LogLevel = "ERROR" // Error level logging.
	FATAL config.LogLevel = "FATAL" // Fatal level logging, which causes application termination.
)

// OmniLogger is the main structure for the logger, holding configuration, context, and drivers.
type OmniLogger struct {
	config  config.Config      // Logger configuration.
	context *model.Context     // Context information for logging.
	drivers []pkg.LoggerDriver // List of logging drivers.
}

// captureStackTrace captures the stack trace of the caller.
func (l *OmniLogger) captureStackTrace() string {
	pc, file, line, ok := runtime.Caller(4) // Skip 4 frames to get the caller's caller.
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", file, line, funcName)
}

// logWritter writes a log message to all configured drivers.
func (l *OmniLogger) logWritter(level config.LogLevel, message string) {
	if enabled, ok := l.config.LogLevels[level]; !ok || !enabled {
		return
	}

	stack := l.captureStackTrace()
	timestamp := time.Now().Format(time.RFC3339)
	messageData := model.MessageData{
		Level:      l.levelToString(level),
		Message:    message,
		StackTrace: stack,
		Context:    l.context,
		Timestamp:  timestamp,
	}

	var wg sync.WaitGroup

	// Write log messages concurrently to all drivers.
	for _, driver := range l.drivers {
		wg.Add(1)

		go func(driver pkg.LoggerDriver) {
			defer wg.Done()
			formattedMessage, err := driver.FormatLog(messageData)
			if err != nil {
				fmt.Println("Error formatting log:", err)
			}
			err = driver.WriteLog(formattedMessage)
			if err != nil {
				fmt.Println("Error writing log:", err)
			}
		}(driver)
	}

	wg.Wait() // Wait for all log writes to complete.
}

func (l *OmniLogger) levelToString(level config.LogLevel) string {
	return string(level)
}

func (l *OmniLogger) AddCustomLogLevel(level config.LogLevel, enabled bool) {
	l.config.LogLevels[level] = enabled
}

func (l *OmniLogger) Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(DEBUG, message)
}

func (l *OmniLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(INFO, message)
}

func (l *OmniLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(WARN, message)
}

func (l *OmniLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(ERROR, message)
}

func (l *OmniLogger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(FATAL, message)
	os.Exit(1)
}

func (l *OmniLogger) Logf(level config.LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(level, message)
	os.Exit(1)
}

func (l *OmniLogger) Log(level config.LogLevel, message string) {
	l.logWritter(level, message)
}

func (l *OmniLogger) Debug(message string) {
	l.logWritter(DEBUG, message)
}

func (l *OmniLogger) Info(message string) {
	l.logWritter(INFO, message)
}

func (l *OmniLogger) Warn(message string) {
	l.logWritter(WARN, message)
}

func (l *OmniLogger) Error(message string) {
	l.logWritter(ERROR, message)
}

func (l *OmniLogger) Fatal(message string) {
	l.logWritter(FATAL, message)
	os.Exit(1)
}
