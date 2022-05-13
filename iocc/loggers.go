package iocc

import (
	"github.com/zbh255/bilog"
	"io"
	"os"
)

var (
	loggers = make(map[string]bilog.Logger, 3)
)

func RegisterAccessLog(writer io.Writer) {
	loggers["accessLog"] = bilog.NewLogger(writer,bilog.DEBUG)
}

func RegisterErrorLog(writer io.Writer) {
	loggers["errorLog"] = bilog.NewLogger(writer,bilog.ERROR,bilog.WithCaller())
}

func GetAccessLog() bilog.Logger {
	return loggers["accessLog"]
}

func GetStdLog() bilog.Logger {
	return loggers["stdLog"]
}

func GetErrorLog() bilog.Logger {
	return loggers["errorLog"]
}

func init() {
	loggers["stdLog"] = bilog.NewLogger(os.Stdout,bilog.PANIC,bilog.WithCaller())
}
