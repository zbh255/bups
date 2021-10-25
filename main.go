package main

import (
	"fmt"
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"os"
)

const PluginPath = "./plugins"
const ConfigPath = "./config/config.toml"

func main() {
	// 处理错误
	defer func() {
		if err := recover();err != nil {
			fmt.Print(err)
		}
	}()
	// 注册插件
	ctx := plugin.NewContext()
	ctx.StdOut = os.Stdout
	// 提供配置文件
	cfg,err := os.OpenFile(ConfigPath,os.O_WRONLY|os.O_RDONLY|os.O_SYNC,0777)
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
	if fn,ok := fnTable["backup.so"];ok {
		fn(PluginPath + "/backup.so")
		delete(fnTable,"backup.so")
	}
	// 注册其他的插件
	for k, v := range fnTable {
		v(PluginPath + "/" + k)
	}
	// 启动初始化插件
	ctx.SetState(plugin.Init)
	// 并发处理参数，如果有插件需要，则交给该插件
}

