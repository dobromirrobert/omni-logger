package pkg

import "omnilogger/model"

type LoggerDriver interface {
	WriteLog(message string) error
	FormatLog(messageData model.MessageData) (string, error)
}
