package logger

import (
	"io"
	"log"
)

/*
	提供日志支持
*/

type Logger interface {
	Debug(message string)
	Info(message string)
	Error(message string)
	Trace(message string)
	Panic(err error)
}

func New(write io.Writer, prefix string) Logger {
	return &LoggerImpl{
		logger: log.New(write, "Error", log.LstdFlags),
		prefix: prefix,
	}
}

type LoggerImpl struct {
	logger *log.Logger
	prefix string
}

func (l *LoggerImpl) Debug(message string) {
	l.logger.SetPrefix(l.prefix + ".Debug:")
	l.logger.Print(message)
}

func (l *LoggerImpl) Info(message string) {
	l.logger.SetPrefix(l.prefix + ".Info:")
	l.logger.Print(message)
}

func (l *LoggerImpl) Error(message string) {
	l.logger.SetPrefix(l.prefix + ".Error:")
	l.logger.Print(message)
}

func (l *LoggerImpl) Trace(message string) {
	l.logger.SetPrefix(l.prefix + ".Trace:")
	l.logger.Print(message)
}

func (l *LoggerImpl) Panic(err error) {
	l.logger.SetPrefix(l.prefix + ".Panic:")
	l.logger.Panic(err)
}
