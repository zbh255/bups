// +build linux darwin windows

package main

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/iocc"
	"github.com/abingzo/bups/plugins/backup"
	"github.com/abingzo/bups/plugins/daemon"
	"github.com/abingzo/bups/plugins/encrypt"
	"github.com/abingzo/bups/plugins/recovery"
	"github.com/abingzo/bups/plugins/upload"
	"github.com/abingzo/bups/plugins/web_config"
	"os"
)

// LoaderPlugin 根据目录加载目录下的所有插件
func LoaderPlugin(configFilePath string) *plugin.Context {
	// 注册插件
	ctx := plugin.NewContext()

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
	hashTable := make(map[string]struct{},len(mainConfig.Project.Install))
	for _,v := range mainConfig.Project.Install {
		hashTable[v] = struct{}{}
	}
	// 注册插件
	iocc.RegisterPlugin(backup.New)
	iocc.RegisterPlugin(daemon.New)
	iocc.RegisterPlugin(encrypt.New)
	iocc.RegisterPlugin(recovery.New)
	iocc.RegisterPlugin(upload.New)
	iocc.RegisterPlugin(web_config.New)
	// 加载插件
	for _,v := range iocc.GetPluginList() {
		tmpPlg := v()
		_,ok := hashTable[tmpPlg.GetName()]
		if ok {
			ctx.RegisterRaw(tmpPlg)
		}
	}
	return ctx
}
