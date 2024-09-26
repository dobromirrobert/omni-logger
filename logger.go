package omnilogger

import (
	"omnilogger/config"
	"omnilogger/model"
	pkg "omnilogger/pkg"
)

var (
	instance *OmniLogger // Singleton instance of OmniLogger.
)

// NewOmniLogger creates and returns a new logger instance with the given configuration, context, and drivers.
func NewOmniLogger(config config.Config, ctx *model.Context, drivers ...pkg.LoggerDriver) *OmniLogger {
	return &OmniLogger{
		config:  config,
		drivers: drivers,
		context: ctx,
	}
}

// AddConfig updates the configuration of the singleton logger instance.
func AddConfig(config config.Config) {
	ensureInstance()
	instance.config = config
}

// AddDriver appends one or more logging drivers to the singleton logger instance.
func AddDriver(drivers ...pkg.LoggerDriver) {
	ensureInstance()
	instance.drivers = append(instance.drivers, drivers...)
}

// GetOmniLoggerWithContext retrieves a copy of the global logger instance with a specified context.
func GetOmniLoggerWithContext(ctx model.Context) (*OmniLogger, error) {
	ensureInstance()
	return &OmniLogger{
		config:  instance.config,
		drivers: instance.drivers,
		context: &ctx,
	}, nil
}

func ensureInstance() {
	if instance == nil {
		instance = &OmniLogger{}
	}
}

func Warnf(format string, args ...interface{}) {
	ensureInstance()
	instance.Warnf(format, args...)
}

// Logf formats and logs a message at a specified log level ussualy used for custome logs
func Logf(level config.LogLevel, format string, args ...interface{}) {
	ensureInstance()
	instance.Logf(level, format, args...)
}

func Debugf(format string, args ...interface{}) {
	ensureInstance()
	instance.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	ensureInstance()
	instance.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	ensureInstance()
	instance.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	ensureInstance()
	instance.Fatalf(format, args...)
}

// Log logs a message at a specified log level usually used for custome logs
func Log(level config.LogLevel, message string) {
	ensureInstance()
	instance.Log(level, message)
}

func Debug(message string) {
	ensureInstance()
	instance.Debug(message)
}

func Info(message string) {
	ensureInstance()
	instance.Info(message)
}

func Warn(message string) {
	ensureInstance()
	instance.Warn(message)
}

func Error(message string) {
	ensureInstance()
	instance.Error(message)
}

func Fatal(message string) {
	ensureInstance()
	instance.Fatal(message)
}
