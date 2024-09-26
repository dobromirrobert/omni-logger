package pkg

import (
	"encoding/json"
	"fmt"
	"omnilogger/model"
	"os"
)

type FileDriver struct {
	file *os.File
}

func NewFileDriver(filePath string) (*FileDriver, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileDriver{file: file}, nil
}

func (d *FileDriver) WriteLog(message string) error {
	_, err := fmt.Fprintln(d.file, message)
	return err
}

func (d *FileDriver) FormatLog(messageData model.MessageData) (string, error) {
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

func (d *FileDriver) Close() {
	d.file.Close()
}
