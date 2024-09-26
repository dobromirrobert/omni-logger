package main

import (
	"omnilogger/omnilogger"
)

const (
	TRACE omnilogger.LogLevel = "TRACE"
)

func main() {
	// Setup log levels configuration
	config := omnilogger.Config{
		LogLevels: map[omnilogger.LogLevel]bool{
			omnilogger.DEBUG: true,
			omnilogger.INFO:  true,
			omnilogger.WARN:  true,
			omnilogger.ERROR: true,
			omnilogger.FATAL: true,
		},
	}

	// Create a new CLI driver
	cliDriver := &omnilogger.CLIDriver{}

	// Create a new File driver that writes in app.log
	fileDriver, err := omnilogger.NewFileDriver("app.log")
	if err != nil {
		panic("Failed to create driver file")
	}

	// initialize the singleton
	omnilogger.AddDriver(cliDriver, fileDriver)
	omnilogger.AddConfig(config)

	// Create a new OmniLogger with the CLI driver
	logger := omnilogger.NewOmniLogger(config, nil, cliDriver)

	// Logging a simple message without context
	logger.Debug("Debug message: System starting")
	logger.Info("Info message: System initialized")
	logger.Warn("Warn message: Low disk space")
	logger.Error("Error message: File not found")

	// Logging a message with context
	context := &omnilogger.Context{
		TransactionID: "tx123",
		UserID:        "user456",
		MetaData: map[string]interface{}{
			"ip":        "192.168.1.1",
			"operation": "file upload",
		},
	}

	// Create a new logger with the context that inherit the singleton drivers and config
	contextLogger, err := omnilogger.GetOmniLoggerWithContext(*context)
	if err != nil {
		panic("Failed to create logger with context")
	}

	// Log a message with the context
	contextLogger.Infof("Info message: Upload started with file %s", "execute.go")
	contextLogger.Errorf("Error message: Upload %s failed due to timeout", "system.go")
	omnilogger.Errorf("Error message: Upload %s failed due to timeout", "system.go")

	//Add custom log
	contextLogger.AddCustomLogLevel(TRACE, true)
	contextLogger.Log(TRACE, "This is a custom log")

}
