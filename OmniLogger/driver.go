package omnilogger

type LoggerDriver interface {
	WriteLog(message string) error
	FormatLog(messageData MessageData) (string, error)
}
