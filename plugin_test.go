package main

import (
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/iocc"
	"os"
	"testing"
)

type TestPlugin struct {
	name    string
	support []uint32
	_type   plugin.Type
}

func (t *TestPlugin) Start(args []string) {
	return
}

func (t *TestPlugin) Caller(single plugin.Single) {
	return
}

func (t *TestPlugin) GetName() string {
	return t.name
}

func (t *TestPlugin) GetType() plugin.Type {
	return t._type
}

func (t *TestPlugin) GetSupport() []uint32 {
	return t.support
}

func (t *TestPlugin) SetSource(source *plugin.Source) {
	return
}

func testRegisterSource() {
	// 注册主配置文件
	configFile, err := os.Open("./config/dev/config.toml")
	if err != nil {
		panic(err)
	}
	iocc.RegisterConfig(configFile)
	config := iocc.GetConfig()
	// 注册日志器
	stdLog := iocc.GetStdLog()
	accessLogFd, err := os.OpenFile(config.Project.Log.AccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		stdLog.ErrorFromString(err.Error())
	}
	errorLogFd, err := os.OpenFile(config.Project.Log.ErrorLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		stdLog.ErrorFromString(err.Error())
	}
	iocc.RegisterAccessLog(accessLogFd)
	iocc.RegisterErrorLog(errorLogFd)
}

func TestPluginLoad(t *testing.T) {
	testRegisterSource()

	// 插件管理的Context
	ctx := plugin.NewContext()
	// 创建原始资源
	rawSource := new(plugin.Source)
	rawSource.StdLog = iocc.GetStdLog()
	rawSource.AccessLog = iocc.GetAccessLog()
	rawSource.ErrorLog = iocc.GetErrorLog()
	file, err := os.OpenFile("./config/dev/config.toml", os.O_RDWR, 0755)
	if err != nil {
		t.Error(err)
		return
	}
	rawSource.RawConfig = NewCFGBuffer(file)
	rawSource.Config = iocc.GetConfig()
	ctx.RawSource = rawSource
	// 注册插件列表
	// 读取配置文件，决定加载那些插件
	mainConfig := iocc.GetConfig()
	hashTable := make(map[string]struct{}, len(mainConfig.Project.Install))
	for _, v := range mainConfig.Project.Install {
		hashTable[v] = struct{}{}
	}
	supportTable := []uint32{
		plugin.SUPPORT_LOGGER,
		plugin.SUPPORT_CONFIG_OBJ,
	}
	// 注册插件
	iocc.RegisterPlugin(func() plugin.Plugin {
		return &TestPlugin{
			name:    "backup",
			support: supportTable,
			_type:   plugin.BHandle,
		}
	})
	iocc.RegisterPlugin(func() plugin.Plugin {
		return &TestPlugin{
			name:    "daemon",
			support: append(supportTable, []uint32{plugin.SUPPORT_ARGS}...),
			_type:   plugin.Init,
		}
	})
	iocc.RegisterPlugin(func() plugin.Plugin {
		return &TestPlugin{
			name:    "encrypt",
			support: supportTable,
			_type:   plugin.BHandle,
		}
	})
	iocc.RegisterPlugin(func() plugin.Plugin {
		return &TestPlugin{
			name:    "upload",
			support: supportTable,
			_type:   plugin.BCallBack,
		}
	})
	iocc.RegisterPlugin(func() plugin.Plugin {
		return &TestPlugin{
			name:    "web_config",
			support: supportTable,
			_type:   plugin.Init,
		}
	})
	// 加载插件
	for _, v := range iocc.GetPluginList() {
		tmpPlg := v()
		_, ok := hashTable[tmpPlg.GetName()]
		if ok {
			ctx.RegisterRaw(tmpPlg)
		}
	}

	// 测试插件管理Context的状态切换
	ctx.RangeAllPlugin(func(k int, v plugin.Plugin) {
		if _, ok := hashTable[v.GetName()]; !ok {
			t.Error("plugin loader name is not equal")
			return
		}
	})
	// 测试所有支持参数的插件是否注册到位
	argsPlugins := make(map[string]struct{}, 5)
	for _, v := range iocc.GetPluginList() {
		for _, supportType := range v().GetSupport() {
			if supportType == plugin.SUPPORT_ARGS {
				argsPlugins[v().GetName()] = struct{}{}
			}
		}
	}
	ctx.RangeArgsPlugin(func(k int, v plugin.Plugin) {
		if _, ok := argsPlugins[v.GetName()]; !ok {
			t.Error("plugin register failed")
			return
		}
	})
	// 测试插件管理Context的状态切换
	stateTable := []plugin.Type{
		plugin.Init,
		plugin.BCollect,
		plugin.BHandle,
		plugin.BCallBack,
	}
	for _, v := range stateTable {
		ctx.SetState(v)
		if ctx.GetState() != v {
			t.Error("test plugin Context state switch failed")
		}
	}
}

// 测试真实的plugin组件加载，并会调用它们的一些接口
func TestTruePluginLoad(t *testing.T) {
	err := os.Chdir("./test")
	if err != nil {
		t.Error(err.Error())
		return
	}
	RegisterSource()
	LoaderPlugin()
}