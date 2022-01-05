package iocc

import (
	"github.com/abingzo/bups/common/plugin"
)

// 存储插件的初始化方法
var pluginList []func() plugin.Plugin


func RegisterPlugin(fn func() plugin.Plugin) {
	pluginList = append(pluginList,fn)
}

func GetPluginList() []func() plugin.Plugin {
	return pluginList
}