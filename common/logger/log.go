package logger

import (
	"github.com/mengzushan/bups/common/error"
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

/*
	关闭一个不存在的文件会引发panic,所以必须使用recover
*/
func (l *Logger) Close() {
	defer func() {
		if err := recover(); err != nil {
			println("There was an error closing the file")
		}
	}()
	err := l.file.Close()
	if err != nil {
		panic(err)
	}
}

/*
	参数必须为string
	参数类型设置为interface{}只是方便传nil值
*/
func Std(p interface{}) (*Logger, error.Error) {
	var path string
	if p != nil {
		path = p.(string)
	} else {
		pathHead, _ := os.Getwd()
		path = pathHead + filepath.FromSlash("/log/app.log")
	}
	//file, _ := os.Open(pathHead + filepath.FromSlash("/log/app.log"))
	// 以追加写入模式打开文件
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	// 日志如果遇到异常则整个程序都会停止
	// 文件不存在则创建文件
	if err != nil {
		_, err = os.Create(path)
		if err != nil {
			return nil, error.LogError + error.Error(err.Error())
		}
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, error.LogError + error.Error(err.Error())
		}
	}
	l := log.New(file, "", log.Lshortfile)
	l.SetFlags(log.LstdFlags)
	return &Logger{
		file:   file,
		logger: l,
	}, error.Nil
}
