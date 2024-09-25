# OmniLogger

OmniLogger is a flexible and powerful logging library for Go, designed to support multiple logging drivers (e.g., console, file) and customizable log levels. It provides a simple interface for logging messages with context information, making it easy to integrate into any application.

## Features

- **Multiple Drivers**: Support for various logging drivers, allowing logs to be written to different outputs (e.g., files, console).
- **Configurable Log Levels**: Easily configure which log levels are enabled or disabled.
- **Contextual Logging**: Attach metadata such as transaction IDs and user IDs to log messages for better traceability.
- **Stack Trace Capture**: Automatically capture and log stack traces for error messages.
