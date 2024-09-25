package omnilogger

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

type MockDriver struct {
	messages []string
}

func (m *MockDriver) WriteLog(message string) error {
	m.messages = append(m.messages, message)
	return nil
}

func (m *MockDriver) FormatLog(messageData MessageData) (string, error) {
	return messageData.Message, nil
}

func TestFileDriver(t *testing.T) {
	// Create a temporary file to write logs
	file, err := os.CreateTemp("", "logfile*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create a new FileDriver
	driver, err := NewFileDriver(file.Name())
	if err != nil {
		t.Fatalf("could not create FileDriver: %v", err)
	}
	defer driver.Close()

	// Write a log message using the FileDriver
	err = driver.WriteLog("This is a test log message")
	if err != nil {
		t.Fatalf("WriteLog failed: %v", err)
	}

	// Check if the message was written to the file
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("could not read temp file: %v", err)
	}

	expected := "This is a test log message\n"

	if string(content) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(content))
	}
}

func TestFileDriverWithFormatting(t *testing.T) {
	// Create a temporary file to write logs
	file, err := os.CreateTemp("", "logfile*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create a new FileDriver
	driver, err := NewFileDriver(file.Name())
	if err != nil {
		t.Fatalf("could not create FileDriver: %v", err)
	}
	defer driver.Close()

	// Create a mock log message data
	context := &Context{
		TransactionID: "tx123",
		UserID:        "user456",
		MetaData: map[string]interface{}{
			"extra": "metadata",
		},
	}

	timestamp := time.Now().Format(time.RFC3339)
	messageData := MessageData{
		Level:      string(INFO),
		Message:    "This is a formatted log message",
		StackTrace: "example.go:123 main.someFunction",
		Context:    context,
		Timestamp:  timestamp,
	}

	// Use FileDriver's FormatLog function to format the log
	formattedLog, err := driver.FormatLog(messageData)
	if err != nil {
		t.Fatalf("FormatLog failed: %v", err)
	}

	// Write the formatted log to the file
	err = driver.WriteLog(formattedLog)
	if err != nil {
		t.Fatalf("WriteLog failed: %v", err)
	}

	// Check if the formatted log was written to the file
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("could not read temp file: %v", err)
	}

	// Expected formatted log as a JSON string
	var logEntry map[string]interface{}
	if err := json.Unmarshal(content, &logEntry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	// Verify log fields
	if logEntry["level"] != string(INFO) {
		t.Errorf("expected level to be 'INFO', got '%v'", logEntry["level"])
	}

	if logEntry["message"] != "This is a formatted log message" {
		t.Errorf("expected message to be 'This is a formatted log message', got '%v'", logEntry["message"])
	}

	if logEntry["stack_trace"] != "example.go:123 main.someFunction" {
		t.Errorf("expected stack_trace to be 'example.go:123 main.someFunction', got '%v'", logEntry["stack_trace"])
	}

	if logEntry["transaction_id"] != "tx123" {
		t.Errorf("expected transaction_id to be 'tx123', got '%v'", logEntry["transaction_id"])
	}

	if logEntry["user_id"] != "user456" {
		t.Errorf("expected user_id to be 'user456', got '%v'", logEntry["user_id"])
	}

	if logEntry["extra"] != "metadata" {
		t.Errorf("expected extra metadata to be 'metadata', got '%v'", logEntry["extra"])
	}

	if logEntry["timestamp"] != timestamp {
		t.Errorf("expected timestamp to be '%v', got '%v'", timestamp, logEntry["timestamp"])
	}
}
func TestOmniLogger(t *testing.T) {
	// Create a mock log driver
	mockDriver := &MockDriver{}

	// Set up config with DEBUG and ERROR levels enabled
	config := Config{
		LogLevels: map[LogLevel]bool{
			DEBUG: true,
			INFO:  false,
			WARN:  false,
			ERROR: true,
		},
	}

	// Create a new OmniLogger with the mock driver
	logger := NewOmniLogger(config, nil, mockDriver)

	// Log a debug message
	logger.Debug("This is a debug message")

	// Log an info message
	logger.Info("This is an info message")

	// Log an error message
	logger.Error("This is an error message")

	// Verify that only DEBUG and ERROR messages were logged
	if len(mockDriver.messages) != 2 {
		t.Fatalf("expected 2 log messages, got %d", len(mockDriver.messages))
	}

	if mockDriver.messages[0] != "This is a debug message" {
		t.Errorf("expected first message to be 'This is a debug message', got %s", mockDriver.messages[0])
	}

	if mockDriver.messages[1] != "This is an error message" {
		t.Errorf("expected second message to be 'This is an error message', got %s", mockDriver.messages[1])
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary JSON config file for testing
	configJSON := `{
		"log_levels": {
			"DEBUG": true,
			"INFO": false,
			"WARN": true,
			"ERROR": true
		}
	}`
	file, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write the test config to the file
	_, err = file.WriteString(configJSON)
	if err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}
	file.Close()

	// Test the LoadConfig function
	config, err := LoadConfig(file.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check if the config was parsed correctly
	if !config.LogLevels[DEBUG] {
		t.Error("expected DEBUG to be enabled")
	}
	if config.LogLevels[INFO] {
		t.Error("expected INFO to be disabled")
	}
	if !config.LogLevels[WARN] {
		t.Error("expected WARN to be enabled")
	}
	if !config.LogLevels[ERROR] {
		t.Error("expected ERROR to be enabled")
	}
}

func TestCLIDriver_FormatLog_NoContext(t *testing.T) {
	// Create a new CLI driver
	driver := &CLIDriver{}

	// Create log message data without a context
	timestamp := time.Now().Format(time.RFC3339)
	messageData := MessageData{
		Level:      string(INFO),
		Message:    "Test message without context",
		StackTrace: "file.go:123",
		Timestamp:  timestamp,
	}

	// Format the log
	formattedLog, err := driver.FormatLog(messageData)
	if err != nil {
		t.Fatalf("FormatLog failed: %v", err)
	}

	// Check if the formatted log contains the correct data
	expectedLog := "[" + string(INFO) + "] timestamp: " + timestamp + "  trace: file.go:123  msg : Test message without context "
	if formattedLog != expectedLog {
		t.Errorf("expected formatted log to be '%s', got '%s'", expectedLog, formattedLog)
	}
}

func TestCLIDriver_FormatLog_WithContext(t *testing.T) {
	// Create a new CLI driver
	driver := &CLIDriver{}

	// Create log message data with context
	context := &Context{
		TransactionID: "tx123",
		UserID:        "user456",
		MetaData: map[string]interface{}{
			"key": "value",
		},
	}
	timestamp := time.Now().Format(time.RFC3339)
	messageData := MessageData{
		Level:      string(INFO),
		Message:    "Test message with context",
		StackTrace: "file.go:123",
		Context:    context,
		Timestamp:  timestamp,
	}

	// Format the log
	formattedLog, err := driver.FormatLog(messageData)
	if err != nil {
		t.Fatalf("FormatLog failed: %v", err)
	}

	// Check if the formatted log contains the correct data
	expectedLog := "[" + string(INFO) + "] timestamp: " + timestamp + " transaction_id: tx123 user_id: user456 key: value  trace: file.go:123  msg : Test message with context "
	if formattedLog != expectedLog {
		t.Errorf("expected formatted log to be '%s', got '%s'", expectedLog, formattedLog)
	}
}

func TestCLIDriver_WriteLog(t *testing.T) {
	// Create a new CLI driver
	driver := &CLIDriver{}

	// Create log message data without a context
	timestamp := time.Now().Format(time.RFC3339)
	messageData := MessageData{
		Level:      string(INFO),
		Message:    "Test CLI message",
		StackTrace: "file.go:123",
		Timestamp:  timestamp,
	}

	// Format the log
	formattedLog, err := driver.FormatLog(messageData)
	if err != nil {
		t.Fatalf("FormatLog failed: %v", err)
	}

	err = driver.WriteLog(formattedLog)
	if err != nil {
		t.Fatalf("WriteLog failed: %v", err)
	}

}
