package main

import (
	"fmt"
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/path"
	"github.com/abingzo/bups/common/plugin"
	"os"
	"time"
)

func main() {
	// 处理错误
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()

	ctx := LoadPlugin(path.PathPluginFileFolder, path.AppLogFilePath, path.PathConfigFile)
	mainConf := config.Read(ctx.Conf).Main
	// 处理参数，如果有插件需要，则交给该插件
	if ArgsProcess(ctx) {
		return
	}
	// 没有参数处理的情况下则通过调度器直接启动程序
	// 启动初始化插件
	ctx.SetState(plugin.Init)
	for {
		timer := time.After(time.Duration(mainConf.LoppTime) * time.Hour)
		select {
		case <-timer:
			ctx.SetState(plugin.BCollect)
			ctx.SetState(plugin.BHandle)
			ctx.SetState(plugin.BCallBack)
		}
	}
}

// 创建文件失败则panic
func createAppLogFile(logFilePath string) *os.File {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_SYNC|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	return file
}
