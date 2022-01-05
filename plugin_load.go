// +build linux darwin windows

package main

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/ioc"
	"github.com/abingzo/bups/plugins/backup"
	"github.com/abingzo/bups/plugins/daemon"
	"github.com/abingzo/bups/plugins/encrypt"
	"github.com/abingzo/bups/plugins/recovery"
	"github.com/abingzo/bups/plugins/upload"
	"github.com/abingzo/bups/plugins/web_config"
	"os"
)

// LoaderPlugin 根据目录加载目录下的所有插件
func LoaderPlugin(logFilePath, configFilePath string) *plugin.Context {
	// 注册插件
	ctx := plugin.NewContext()
	var ReadLogFile = func() *os.File {
		// 准备日志文件
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			panic(err)
		}
		return file
	}
	ctx.LogOut = ReadLogFile
	ctx.StdOut = os.Stdout
	// 提供配置文件
	cfg, err := os.OpenFile(configFilePath, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer cfg.Close()
	// 初始化可以多次读写的配置文件接口
	cfgE := &CFG{}
	cfgE.Open(cfg)
	ctx.Conf = cfgE
	// 读取配置文件，决定加载那些插件
	mainConfig  := config.Read(ctx.Conf)
	hashTable := make(map[string]struct{},len(mainConfig.Main.Install))
	for _,v := range mainConfig.Main.Install {
		hashTable[v] = struct{}{}
	}
	// 注册插件
	ioc.RegisterPlugin(backup.New)
	ioc.RegisterPlugin(daemon.New)
	ioc.RegisterPlugin(encrypt.New)
	ioc.RegisterPlugin(recovery.New)
	ioc.RegisterPlugin(upload.New)
	ioc.RegisterPlugin(web_config.New)
	// 加载插件
	for _,v := range ioc.GetPluginList() {
		tmpPlg := v()
		_,ok := hashTable[tmpPlg.GetName()]
		if ok {
			ctx.RegisterRaw(tmpPlg)
		}
	}
	return ctx
}
