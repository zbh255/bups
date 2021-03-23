package logger

import (
	"fmt"
	"github.com/mengzushan/bups/utils"
	"log"
	"os"
)

type STD interface {
	StdDebugLog(info string)
	StdInfoLog(info string)
	StdErrorLog(info string)
	StdWarnLog(info string)
	StdTraceLog(info string)
	Close()
}

type Logger struct {
	file 	*os.File
	logger 	*log.Logger
}

func (l *Logger)StdDebugLog(info string)  {
	l.logger.SetPrefix("[Debug]")
	l.logger.Println(info)
}

func (l *Logger)StdInfoLog(info string) {
	l.logger.SetPrefix("[Info]")
	l.logger.Println(info)
}

func (l *Logger)StdErrorLog(info string) {
	l.logger.SetPrefix("[Error]")
	l.logger.Fatal(info)
}

func (l *Logger)StdWarnLog(info string)  {
	l.logger.SetPrefix("[Warn]")
	l.logger.Println(info)
}

func (l *Logger)StdTraceLog(info string) {
	l.logger.SetPrefix("[Trace]")
	l.logger.Println(info)
}

func (l *Logger)Close() {
	err := l.file.Close()
	if err != nil {
		panic(err)
	}
}

func Std()  STD {
	pathHead,_ := os.Getwd()
	pathRuned := utils.PathRune()
	file, _ := os.Open(pathHead + fmt.Sprintf("%slog%sapp.log", pathRuned, pathRuned))
	l := log.New(file,"",log.Lshortfile)
	l.SetFlags(log.LstdFlags)
	var s STD = &Logger{
		file: file,
		logger: l,
	}
	return s
}

