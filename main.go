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

	// 注册插件
	ctx := plugin.NewContext()
	// 准备日志文件
	file, err := os.OpenFile(path.AppLogFilePath, os.O_APPEND|os.O_SYNC, 0777)
	if err != nil {
		file = createAppLogFile()
	}
	ctx.LogOut = file
	// 提供配置文件
	cfg, err := os.OpenFile(path.PathConfigFile, os.O_RDWR|os.O_SYNC, 0777)
	if err != nil {
		panic(err)
	}
	ctx.Conf = cfg
	// 注册在配置文件中声明的插件
	mainConf := config.Read(ctx.Conf).Main
	fnTable := make(map[string]func(string))
	for _, v := range mainConf.Install {
		fnTable[v] = ctx.Register
	}
	// 优先注册备份插件
	if fn, ok := fnTable["backup.so"]; ok {
		fn(path.PathPluginFileFolder + "/backup.so")
		delete(fnTable, "backup.so")
	}
	// 注册其他的插件
	for k, v := range fnTable {
		v(path.PathPluginFileFolder + "/" + k)
	}
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
			ctx.SetState(plugin.BCallBack)
		}
	}
}

// 创建文件失败则panic
func createAppLogFile() *os.File {
	file, err := os.OpenFile(path.AppLogFilePath, os.O_CREATE|os.O_SYNC|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	return file
}
