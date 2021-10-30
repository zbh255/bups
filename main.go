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
	ctx.LogOut = os.Stdout
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
			ctx.SetState(plugin.BStart)
			ctx.SetState(plugin.BCallBack)
		}
	}
}
