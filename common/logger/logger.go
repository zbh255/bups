package logger

import (
	"io"
	"log"
)

// 日志等级的定义
const (
	INFO uint8 = iota
	DEBUG
	TRACE
	ERROR
	PANIC
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

func New(write io.Writer, level uint8) Logger {
	return &LoggerImpl{
		logger: log.New(write, "Error", log.LstdFlags | log.Lshortfile),
		level:       level,
	}
}

type LoggerImpl struct {
	logger *log.Logger
	level uint8
}

func (l *LoggerImpl) checkLevel(level uint8) bool {
	return l.level >= level
}

func (l *LoggerImpl) Debug(message string) {
	if !l.checkLevel(DEBUG) {
		return
	}
	l.logger.SetPrefix("[Debug]")
	l.logger.Print(message)
}

func (l *LoggerImpl) Info(message string) {
	if !l.checkLevel(INFO) {
		return
	}
	l.logger.SetPrefix("[Info]")
	l.logger.Print(message)
}

func (l *LoggerImpl) Error(message string) {
	if !l.checkLevel(ERROR) {
		return
	}
	l.logger.SetPrefix("[Error]")
	l.logger.Print(message)
}

func (l *LoggerImpl) Trace(message string) {
	if !l.checkLevel(TRACE) {
		return
	}
	l.logger.SetPrefix("[Trace]")
	l.logger.Print(message)
}

func (l *LoggerImpl) Panic(err error) {
	if !l.checkLevel(PANIC) {
		return
	}
	l.logger.SetPrefix("[Panic]")
	l.logger.Panic(err)
}
