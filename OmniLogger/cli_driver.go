package omnilogger

import (
	"fmt"
)

type CLIDriver struct{}

func (d *CLIDriver) WriteLog(message string) error {
	fmt.Println(message)
	return nil
}

func (d *CLIDriver) FormatLog(messageData MessageData) (string, error) {

	logEntry := fmt.Sprintf("[%s] timestamp: %s ", messageData.Level, messageData.Timestamp)

	if messageData.Context == nil {
		logEntry += fmt.Sprintf(" trace: %+s  msg : %s ", messageData.StackTrace, messageData.Message)
		return logEntry, nil
	}

	if messageData.Context.TransactionID != "" {
		logEntry += fmt.Sprintf("transaction_id: %s ", messageData.Context.TransactionID)
	}
	if messageData.Context.UserID != "" {
		logEntry += fmt.Sprintf("user_id: %s ", messageData.Context.UserID)
	}
	for key, value := range messageData.Context.MetaData {
		logEntry += fmt.Sprintf("%s: %v ", key, value)
	}

	logEntry += fmt.Sprintf(" trace: %+s  msg : %s ", messageData.StackTrace, messageData.Message)

	return logEntry, nil
}
