package model

// MessageData represents the structure of a log message.
type MessageData struct {
	Level      string   // The log level of the message.
	Message    string   // The actual log message.
	StackTrace string   // The stack trace at the time of logging.
	Context    *Context // Contextual information for the log entry.
	Timestamp  string   // The timestamp when the log entry was created.
}

// Context holds contextual information for a log entry.
type Context struct {
	TransactionID string                 // Unique identifier for the transaction.
	UserID        string                 // Identifier for the user associated with the log.
	MetaData      map[string]interface{} // Additional metadata related to the log entry.
}
