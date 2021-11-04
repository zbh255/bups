package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/abingzo/bups/common/plugin"
	"os"
	"strings"
)

/*
	处理主程序的参数和插件之间的参数传递的问题
	一个正确的插件参数: ./bups --plugin=daemon --args '<--s start>'
	一个正确的程序参数: ./bups --option pluginInstallList
*/

// ArgsProcess 插件收到的标准参数:
// 原参数:/User/harder/bups --plugin daemon --args '<--s start>'
// 插件看到的:/User/harder/bups --s start
// 有处理参数则返回true,没有处理与调度器无关的程序则返回false
func ArgsProcess(ctx *plugin.Context) bool {
	tag := false
	pluginName := flag.String("plugin", "", "调用的插件的名字")
	caller := flag.String("caller", "", "直接调用一个插件,没有参数传递")
	pluginArgs := flag.String("args", "", "传递的插件参数，比如:'<--s stop>'")
	option := flag.String("option", "", "应用程序选项: pluginInstallList 列出所有安装的插件")
	flag.Parse()
	switch *option {
	case "pluginInstallList":
		tag = true
		ctx.RangeAllPlugin(func(k int, v plugin.Plugin) {
			fmt.Printf("Handler:%d --> PluginName:%s --> PluginType:%d\n", k, v.GetName(), v.GetType())
		})
	case "":
		break
	case "version":
		tag = true
		v := GetInfo()
		bytes,err := json.MarshalIndent(v,"","\t")
		if err != nil {
			fmt.Printf("%s",err.Error())
		} else {
			fmt.Print(string(bytes))
		}
	default:
		tag = true
		break
	}
	// 需要处理传递给插件的参数
	if *pluginName != "" && *pluginArgs != "" {
		// 找到接收的插件
		tag = true
		ctx.RangeArgsPlugin(func(k int, v plugin.Plugin) {
			if v.GetName() == *pluginName {
				v.Start(MainAppArgsToPlugin(*pluginArgs))
			}
		})
	} else if *caller != "" {
		tag = true
		ctx.RangeAllPlugin(func(k int, v plugin.Plugin) {
			if v.GetName() == *caller {
				v.Start(nil)
			}
		})
	}
	return tag
}

// MainAppArgsToPlugin 该函数将主程序参数转换为插件参数
func MainAppArgsToPlugin(s string) []string {
	args := []string{os.Args[0]}
	s = s[1 : len(s)-1]
	args = append(args, strings.Split(s, " ")...)
	return args
}
