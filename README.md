# OmniLogger

OmniLogger is a flexible and powerful logging library for Go, designed to support multiple logging drivers (e.g., console, file) and customizable log levels. It provides a simple interface for logging messages with context information, making it easy to integrate into any application.

## Features

- **Multiple Drivers**: Support for various logging drivers, allowing logs to be written to different outputs (e.g., files, console). Also allow to add new drivers by implementing the interface Driver.
- **Configurable Log Levels**: Easily configure which log levels are enabled or disabled.
- **Contextual Logging**: Attach metadata such as transaction IDs and user IDs to log messages for better traceability.
- **Stack Trace Capture**: Automatically capture and log stack traces for error messages.

## Usage

### Creating a Logger
OmniLogger requires a configuration file or struct that specifies which log levels should be enabled. You can create a Config struct to manage these settings or use LoadConfig function provided to read from a json file. You can create an instance of the logger by passing a configuration and one or more log drivers (e.g., CLI driver for console output).
    ```go
    
    cliDriver := &omnilogger.CLIDriver{} 
    logger := omnilogger.NewOmniLogger(config, nil, cliDriver)

### Logging With Context
OmniLogger also supports logging with additional context (e.g., transaction ID, user ID, and other metadata). Here’s how you can log with context:

  ```go
    context := &omnilogger.Context{
        TransactionID: "tx123",
        UserID:        "user456",
        MetaData: map[string]interface{}{
            "ip": "192.168.1.1",
            "operation": "file upload",
        },
    }

    contextLogger, err := omnilogger.GetOmniLoggerWithContext(*context)
    if err != nil {
        panic("Failed to create logger with context")
    }

    
    contextLogger.Infof("Info message: File uploaded successfully")
    contextLogger.Errorf("Error message: Upload failed")
```


### Custom Log Drivers
You can add custom log drivers by implementing the LoggerDriver interface. For example, you can create a file-based log driver or any other log driver that supports the required WriteLog and FormatLog methods.

```go
type MyCustomDriver struct {}

func (d *MyCustomDriver) WriteLog(message string) error {
    // Write log to custom destination
    return nil
}

func (d *MyCustomDriver) FormatLog(messageData omnilogger.MessageData) (string, error) {
    // Format the log message
    return "formatted message", nil
}
```
### Diagram

![Diagram](https://github.com/user-attachments/assets/d5f532d1-184a-476e-be0c-b9f4f23184ca)
  


