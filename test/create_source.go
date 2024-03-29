package test

import (
	"github.com/abingzo/bups/app"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/iocc"
	"os"
)

// 准备测试包的资源
//go:generate make delete_file
//go:generate make create_file

func testRegisterSource() {
	// 注册主配置文件
	configFile, err := os.Open("./config.toml")
	if err != nil {
		panic(err)
	}
	iocc.RegisterConfig(configFile)
	config := iocc.GetConfig()
	// 注册日志器
	// 测试用的日志器直接使用Stdout&StdErr
	iocc.RegisterAccessLog(os.Stdout)
	iocc.RegisterErrorLog(os.Stderr)
	// 为插件创建自己对应的私有目录
	for _,v := range config.Project.Install {
		err = os.MkdirAll("./cache/"+v, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func LoadPluginSource() *plugin.Source {
	// 创建原始资源
	rawSource := new(plugin.Source)
	rawSource.StdLog = iocc.GetStdLog()
	rawSource.AccessLog = iocc.GetAccessLog()
	rawSource.ErrorLog = iocc.GetErrorLog()
	file, err := os.OpenFile("./config.toml", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	rawSource.RawConfig = app.NewCFGBuffer(file)
	rawSource.Config = iocc.GetConfig()
	return rawSource
}

func init()  {
	testRegisterSource()
}