package app

import (
	"github.com/abingzo/bups/iocc"
	"os"
)

// RegisterSource 负责将所有用到的资源装载进iocc中
func RegisterSource(configPath string) {
	// 注册主配置文件
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	iocc.RegisterConfig(configFile)
	config := iocc.GetConfig()
	// 注册日志器
	stdLog := iocc.GetStdLog()
	accessLogFd, err := os.OpenFile(config.Project.Log.AccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		stdLog.ErrorFromString(err.Error())
	}
	errorLogFd, err := os.OpenFile(config.Project.Log.ErrorLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		stdLog.ErrorFromString(err.Error())
	}
	iocc.RegisterAccessLog(accessLogFd)
	iocc.RegisterErrorLog(errorLogFd)
}
