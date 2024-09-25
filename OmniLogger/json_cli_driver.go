package omnilogger

import (
	"encoding/json"
	"fmt"
)

// JsonCliDriver implements the LoggerDriver interface for JSON output to CLI
type JsonCliDriver struct{}

func (d *JsonCliDriver) WriteLog(message string) error {
	fmt.Println(message)
	return nil
}

func (d *JsonCliDriver) FormatLog(messageData MessageData) (string, error) {
	logEntry := map[string]interface{}{
		"level":     messageData.Level,
		"timestamp": messageData.Timestamp,
	}
	if messageData.Context != nil {
		if messageData.Context.TransactionID != "" {
			logEntry["transaction_id"] = messageData.Context.TransactionID
		}
		if messageData.Context.UserID != "" {
			logEntry["user_id"] = messageData.Context.UserID
		}
		for key, value := range messageData.Context.MetaData {
			logEntry[key] = value
		}
	}
	logEntry["stack_trace"] = messageData.StackTrace
	logEntry["message"] = messageData.Message

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
