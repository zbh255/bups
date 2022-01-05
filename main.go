package main

import (
	"fmt"
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/path"
	"github.com/abingzo/bups/common/plugin"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 处理错误
	defer func() {
		if err := recover(); err != nil {
			stack := stack(3)
			fmt.Printf("PANIC: %s\n%s", err, stack)
		}
	}()

	ctx := LoaderPlugin(path.AppLogFilePath, path.PathConfigFile)
	// 为插件准备存放文件的文件夹，已存在则不创建
	ctx.RangeAllPlugin(func(k int, v plugin.Plugin) {
		info, err := os.Stat(path.PathBackUpCache + "/" + v.GetName())
		if err == nil && info.IsDir() {
			return
		}
		// 不存在则创建
		err = os.MkdirAll(path.PathBackUpCache+"/"+v.GetName(), 0755)
		if err != nil {
			panic(err)
		}
	})
	mainConf := config.Read(ctx.Conf).Main
	// 处理参数，如果有插件需要，则交给该插件
	if ArgsProcess(ctx) {
		return
	}
	// 正常启动的情况下接收信号，并将通知信号派发给插件
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
		switch v := <-c; v {
		case syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT:
			ctx.RangeAllPlugin(func(k int, v plugin.Plugin) {
				v.Caller(plugin.Exit)
			})
			os.Exit(0)
		}
	}()
	// TODO:解决初始化正常却无法打印日志的问题

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
