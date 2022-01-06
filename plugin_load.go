// +build linux darwin windows

package main

import (
	"github.com/abingzo/bups/common/path"
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
// 并为Context初始化原始资源
func LoaderPlugin() *plugin.Context {
	// 注册插件
	ctx := plugin.NewContext()
	// 初始化插件需要的原始资源
	rawSource := new(plugin.Source)
	rawSource.AccessLog = iocc.GetAccessLog()
	rawSource.ErrorLog = iocc.GetErrorLog()
	rawSource.StdLog = iocc.GetStdLog()
	rawSource.Config = iocc.GetConfig()
	// 创建配置文件的原始接口
	configFd, err := os.OpenFile(path.PathConfigFile,os.O_RDWR,0755)
	if err != nil {
		panic(err)
	}
	// 创建一个可重复读取的原始配置文件抽象
	rawSource.RawConfig = NewCFGBuffer(configFd)

	ctx.RawSource = rawSource
	// 读取配置文件，决定加载那些插件
	mainConfig  := iocc.GetConfig()
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
