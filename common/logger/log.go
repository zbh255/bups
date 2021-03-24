package logger

import (
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func (l *Logger) StdDebugLog(info string) {
	l.logger.SetPrefix("[Debug]")
	l.logger.Println(info)
}

func (l *Logger) StdInfoLog(info string) {
	l.logger.SetPrefix("[Info]")
	l.logger.Println(info)
}

func (l *Logger) StdErrorLog(info string) {
	l.logger.SetPrefix("[Error]")
	l.logger.Println(info)
}

func (l *Logger) StdWarnLog(info string) {
	l.logger.SetPrefix("[Warn]")
	l.logger.Println(info)
}

func (l *Logger) StdTraceLog(info string) {
	l.logger.SetPrefix("[Trace]")
	l.logger.Println(info)
}

func (l *Logger) Close() {
	l.file.Close()
}

func Std() *Logger {
	pathHead, _ := os.Getwd()
	file, _ := os.Create(pathHead + filepath.FromSlash("/log/app.log"))
	l := log.New(file, "", log.Lshortfile)
	l.SetFlags(log.LstdFlags)
	return &Logger{
		file: file,
		logger: l,
	}
}
