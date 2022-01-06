package iocc

import (
	"github.com/abingzo/bups/common/logger"
	"io"
	"os"
)

var (
	loggers = make(map[string]logger.Logger,3)
)

func RegisterAccessLog(writer io.Writer) {
	loggers["accessLog"] = logger.New(writer,logger.DEBUG)
}

func RegisterErrorLog(writer io.Writer) {
	loggers["errorLog"] = logger.New(writer,logger.ERROR)
}

func GetAccessLog() logger.Logger {
	return loggers["accessLog"]
}

func GetStdLog() logger.Logger {
	return loggers["stdLog"]
}

func GetErrorLog() logger.Logger {
	return loggers["errorLog"]
}

func init() {
	loggers["stdLog"] = logger.New(os.Stdout, logger.PANIC)
}
