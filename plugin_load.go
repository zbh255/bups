package main

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"os"
)

// LoadPlugin 根据目录加载目录下的所有插件
func LoadPlugin(pluginPh, logFilePath, configFilePath string) *plugin.Context {
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
	// 注册在配置文件中声明的插件
	mainConf := config.Read(ctx.Conf).Main
	fnTable := make(map[string]func(string))
	for _, v := range mainConf.Install {
		fnTable[v] = ctx.Register
	}
	// 优先注册备份插件
	if fn, ok := fnTable["backup.so"]; ok {
		fn(pluginPh + "/backup.so")
		delete(fnTable, "backup.so")
	}
	// 注册其他的插件
	for k, v := range fnTable {
		v(pluginPh + "/" + k)
	}
	return ctx
}
