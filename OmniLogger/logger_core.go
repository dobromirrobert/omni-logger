package omnilogger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// LogLevel represents the level of logging.
type LogLevel string

// Context holds contextual information for a log entry.
type Context struct {
	TransactionID string                 // Unique identifier for the transaction.
	UserID        string                 // Identifier for the user associated with the log.
	MetaData      map[string]interface{} // Additional metadata related to the log entry.
}

const (
	DEBUG LogLevel = "DEBUG" // Debug level logging.
	INFO  LogLevel = "INFO"  // Information level logging.
	WARN  LogLevel = "WARN"  // Warning level logging.
	ERROR LogLevel = "ERROR" // Error level logging.
	FATAL LogLevel = "FATAL" // Fatal level logging, which causes application termination.
)

// MessageData represents the structure of a log message.
type MessageData struct {
	Level      string   // The log level of the message.
	Message    string   // The actual log message.
	StackTrace string   // The stack trace at the time of logging.
	Context    *Context // Contextual information for the log entry.
	Timestamp  string   // The timestamp when the log entry was created.
}

// OmniLogger is the main structure for the logger, holding configuration, context, and drivers.
type OmniLogger struct {
	config  Config         // Logger configuration.
	context *Context       // Context information for logging.
	drivers []LoggerDriver // List of logging drivers.
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
func (l *OmniLogger) logWritter(level LogLevel, message string) {
	if ok := l.config.LogLevels[level]; !ok {
		return // Skip logging if the log level is not enabled.
	}

	stack := l.captureStackTrace()
	timestamp := time.Now().Format(time.RFC3339)
	messageData := MessageData{
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

		go func(driver LoggerDriver) {
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

func (l *OmniLogger) levelToString(level LogLevel) string {
	return string(level)
}

func (l *OmniLogger) AddCustomLogLevel(level LogLevel, enabled bool) {
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

func (l *OmniLogger) Logf(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logWritter(level, message)
	os.Exit(1)
}

func (l *OmniLogger) Log(level LogLevel, message string) {
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
